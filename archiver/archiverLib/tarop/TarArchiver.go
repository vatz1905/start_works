package tarop

import (
	"archive/tar"

	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type FileMetadata struct {
	Fid      int
	FileName string
	Size     int64
	Mode     int64
	ModTime  time.Time
}

var (
	FilesMetadata []FileMetadata
	filesInArch   int64
	tarFile       string
	fullURLFile   string
)

// Create local archive.tar

func CreateTar(files []string, dest string) error {
	// Create and add some files to the archive.
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)

	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			log.Fatal(err)
		}
		fi, err := f.Stat()
		if err != nil {
			log.Fatal(err)
		}
		hdr := &tar.Header{
			Name: file,
			Mode: 0600,
			Size: fi.Size(),
		}
		if err := tw.WriteHeader(hdr); err != nil {
			log.Fatal(err)
		}
		content, err := ioutil.ReadFile(file)
		if err != nil {
			log.Fatal(err)
		}
		if _, err := tw.Write(content); err != nil {
			log.Fatal(err)
		}
	}
	if err := tw.Close(); err != nil {
		log.Fatal(err)
	}
	//Write output to file
	err := ioutil.WriteFile(dest, buf.Bytes(), 0600)
	if err != nil {
		panic(err)
	}

	return nil
}

// Create local archive.targz
func CreateTargz(files []string, dest string) error {
	var buf io.Writer
	// Create output file
	out, err := os.Create(dest)
	if err != nil {
		log.Fatalln("Error writing archive:", err)
	}
	defer out.Close()
	buf = io.Writer(out)
	// Create new Writers for gzip and tar
	// These writers are chained. Writing to the tar writer will
	// write to the gzip writer which in turn will write to
	// the "buf" writer
	gw := gzip.NewWriter(buf)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()
	// Iterate over files and add them to the tar archive
	for _, file := range files {
		err := addToTar(tw, file)
		if err != nil {
			return err
		}
	}
	return nil
}

// Extract tar archives (of types: ".tar" , ".targz", ".tar.gz") of local path archive or url path archive
func Extract(srcArchive string, isUrl bool, destDirPath string) (fileName string, err error) {
	var isTargz bool
	var body []byte
	var tr *tar.Reader

	isTargz = strings.HasSuffix(srcArchive, ".tar.gz") || strings.HasSuffix(srcArchive, ".tgz")

	if !isUrl { //local archive
		archive, err := os.Open(srcArchive)
		if err != nil {
			log.Fatal(err)
		}
		if !isTargz { // local archive.tar
			tr = tar.NewReader(archive)
		} else { // local archive.targz or archive.tar.gz
			archiveTargz, err := gzip.NewReader(archive)
			if err != nil {
				fmt.Println("There is a problem with os.Open")
			}
			tr = tar.NewReader(archiveTargz)
		}
	} else if isUrl { //url archive
		fileURL, err := url.Parse(srcArchive)
		if err != nil {
			log.Fatal(err)
		}
		path := fileURL.Path
		segments := strings.Split(path, "/")
		fileName = segments[len(segments)-1]

		// http get request
		resp, err := http.Get(srcArchive)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		if isTargz { // url archive.targz or url archive.tar.gz
			reader, err := gzip.NewReader(resp.Body)
			if err != nil {
				log.Fatal(err)
			}
			defer reader.Close()
			// wrap reader for tar.NewReader
			body, err = ioutil.ReadAll(reader)
			if err != nil {
				log.Fatal(err)
			}
		} else { // url archive.tar
			// wrap reader for tar.NewReader
			body, err = ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}
		}
		tr = tar.NewReader(bytes.NewReader(body))
		if err != nil {
			log.Fatal(err)
		}
	}

	err = os.Mkdir(destDirPath, 0755)
	if err != nil {
		log.Fatal(err)
	}
	for {

		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		// the target location where the dir/file should be created
		target := filepath.Join(destDirPath, hdr.Name)
		// check the file type
		switch hdr.Typeflag {

		// if its a dir and it doesn't exist create it
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {

				if err := os.MkdirAll(target, 0755); err != nil {
					log.Fatal(err)
				}
			}

		// if it's a file create it
		case tar.TypeReg:

			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(hdr.Mode))
			if err != nil {
				log.Fatal(err)
			}

			// copy over contents
			if _, err := io.Copy(f, tr); err != nil {
				log.Fatal(err)
			}

			f.Close()
		}

	}

	return srcArchive, err
}

// ExtractSome files from tar archives (of types: ".tar" , ".targz", ".tar.gz") of local path archive or url path archive

func ExtractSome(srcArchive string, isUrl bool, targetFiles []string, destDirPath string) (fileName string, err error) {
	var isTargz bool
	var body []byte
	var tr *tar.Reader

	isTargz = strings.HasSuffix(srcArchive, ".tar.gz") || strings.HasSuffix(srcArchive, ".tgz")

	if !isUrl { //local archive
		archive, err := os.Open(srcArchive)
		if err != nil {
			log.Fatal(err)
		}
		if !isTargz { // local archive.tar
			tr = tar.NewReader(archive)
		} else { // local archive.targz or archive.tar.gz
			archiveTargz, err := gzip.NewReader(archive)
			if err != nil {
				fmt.Println("There is a problem with os.Open")
			}
			tr = tar.NewReader(archiveTargz)
		}
	} else if isUrl { //url archive
		fileURL, err := url.Parse(srcArchive)
		if err != nil {
			log.Fatal(err)
		}
		path := fileURL.Path
		segments := strings.Split(path, "/")
		fileName = segments[len(segments)-1]

		// http get request
		resp, err := http.Get(srcArchive)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		if isTargz { // url archive.targz or url archive.tar.gz
			reader, err := gzip.NewReader(resp.Body)
			if err != nil {
				log.Fatal(err)
			}
			defer reader.Close()
			// wrap reader for tar.NewReader
			body, err = ioutil.ReadAll(reader)
			if err != nil {
				log.Fatal(err)
			}
		} else { // url archive.tar
			// wrap reader for tar.NewReader
			body, err = ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}
		}
		tr = tar.NewReader(bytes.NewReader(body))
		if err != nil {
			log.Fatal(err)
		}
	}

	//create destDirPath for extractedFiles
	err = os.Mkdir(destDirPath, 0755)
	if err != nil {
		log.Fatal(err)
	}

	for _, target := range targetFiles {
		for {

			hdr, err := tr.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			//case tar.TypeReg
			if hdr.Typeflag == tar.TypeReg {
				if hdr.Name == target {

					exFile := filepath.Base(target)
					out := filepath.Join(destDirPath, exFile)

					f, err := os.OpenFile(out, os.O_CREATE|os.O_RDWR, os.FileMode(hdr.Mode))
					if err != nil {
						log.Fatal(err)
					}
					// copy over contents
					if _, err := io.Copy(f, tr); err != nil {
						log.Fatal(err)
					}
					f.Close()
					break
				}
			}

		}
	}

	return fileName, err
}

// Get metadata from tar archive (of types: ".tar" , ".targz", ".tar.gz") of local path archive or url path archive
func Meta(srcArchive string, isUrl bool) (tarMetadata []FileMetadata, err error) {
	var isTargz bool
	var body []byte
	var tr *tar.Reader

	filesInArch = 0
	// allocate dynamic array, start with minimum size
	FilesMetadata := make([]FileMetadata, 1)

	isTargz = strings.HasSuffix(srcArchive, ".tar.gz") || strings.HasSuffix(srcArchive, ".tgz")

	if !isUrl { //local archive
		archive, err := os.Open(srcArchive)
		if err != nil {
			log.Fatal(err)
		}
		if !isTargz { // local archive.tar
			tr = tar.NewReader(archive)
		} else { // local archive.targz or archive.tar.gz
			archive_gz, err := gzip.NewReader(archive)
			if err != nil {
				fmt.Println("There is a problem with os.Open")
			}
			tr = tar.NewReader(archive_gz)
		}
	} else if isUrl { //url archive

		// http get request
		resp, err := http.Get(srcArchive)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		if isTargz { // url archive.targz or url archive.tar.gz
			reader, err := gzip.NewReader(resp.Body)
			if err != nil {
				log.Fatal(err)
			}
			defer reader.Close()
			// wrap reader for tar.NewReader
			body, err = ioutil.ReadAll(reader)
			if err != nil {
				log.Fatal(err)
			}
		} else { // url archive.tar
			// wrap reader for tar.NewReader
			body, err = ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}
		}
		tr = tar.NewReader(bytes.NewReader(body))
		if err != nil {
			log.Fatal(err)
		}
	}

	var i int
	i = 0
	for {

		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		filesInArch++
		FilesMetadata[i].Fid = i
		FilesMetadata[i].FileName = hdr.Name
		FilesMetadata[i].Size = hdr.Size
		FilesMetadata[i].Mode = hdr.Mode
		FilesMetadata[i].ModTime = hdr.ModTime

		// resize dynamic array
		FilesMetadata = append(FilesMetadata, FilesMetadata[i])
		i++
	}

	//delete last element in filesMetadata
	FilesMetadata = FilesMetadata[0:filesInArch]
	//prepare output in json format
	b, err := json.Marshal(FilesMetadata)
	if err != nil {
		fmt.Println("error:", err)
	}
	os.Stdout.Write(b)

	return FilesMetadata, nil
}

// helper functions

// help function for CreateTargz function
func addToTar(tw *tar.Writer, fileName string) error {
	// Open the file which will be written into the archive
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	// Get FileInfo about our file providing file size, mode, etc.
	info, err := file.Stat()
	if err != nil {
		return err
	}
	// Create a tar Header from the FileInfo data
	header, err := tar.FileInfoHeader(info, info.Name())
	if err != nil {
		return err
	}
	header.Name = fileName
	// Write file header to the tar archive
	err = tw.WriteHeader(header)
	if err != nil {
		return err
	}
	// Copy file content to tar archive
	_, err = io.Copy(tw, file)
	if err != nil {
		return err
	}
	return nil
}

// download tar archive by its url, to pwd. (for future use)
func GetTarByUrl(tarArchiveUrl string) (fileName string, err error) {

	fileURL, err := url.Parse(tarArchiveUrl)
	if err != nil {
		log.Fatal(err)
	}
	path := fileURL.Path
	segments := strings.Split(path, "/")
	fileName = segments[len(segments)-1]

	// Create blank file
	f, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	// http get request
	resp, err := http.Get(tarArchiveUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if _, err := io.Copy(f, resp.Body); err != nil {
		log.Fatal(err)
	}
	return fileName, err
}
