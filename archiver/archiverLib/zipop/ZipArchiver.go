package zipop

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type FileMetadata struct {
	Fid            int
	FileName       string
	CompressMethod uint16 // Dflate or Store
	CSize          uint64
	UcSize         uint64
	StartInArch    uint64
	EndInArch      uint64
	CRC32          uint32
}

var FilesMetadata []FileMetadata

var (
	filesInArch int64
	zipArchive  string
	dstZipDir   string
	files       []string
)

// create zip archive
func Create(files []string, dest string) error {
	// Create a buffer to write our archive to.
	buf := new(bytes.Buffer)
	// Create a new zip archive.
	w := zip.NewWriter(buf)
	// Add some files to the archive.
	for _, file := range files {
		f, err := w.Create(file)
		if err != nil {
			log.Fatal(err)
		}
		content, err := ioutil.ReadFile(file)
		if err != nil {
			log.Fatal(err)
		}
		_, err = f.Write([]byte(content))
		if err != nil {
			log.Fatal(err)
		}
	}
	// Make sure to check the error on Close.
	err := w.Close()
	if err != nil {
		log.Fatal(err)
	}
	//Write output to file
	err = ioutil.WriteFile(dest, buf.Bytes(), os.ModePerm)
	if err != nil {
		panic(err)
	}
	return err
}

// Extract files from  zip srcArchive , to dest directory
// to do - use dest_dir_path.  now implemented save unzipped archive in folder with same name as src archive trimmed of ".zip"
func Extract(srcArchive string, isUrl bool, destDirPath string) ([]string, error) {

	var rp *zip.Reader
	var filesNames []string

	if !isUrl { //local archive
		r_c, err := zip.OpenReader(srcArchive)
		if err != nil {
			return filesNames, err
		}
		defer r_c.Close()
		r := zip.Reader(r_c.Reader) //temp
		rp = &r
		// mk directory  for unarchive files

	} else if isUrl { //url archive

		err := os.Mkdir(destDirPath, 0755)
		if err != nil {
			log.Fatal(err)
		}

		resp, err := http.Get(srcArchive)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		rp, err = zip.NewReader(bytes.NewReader(body), int64(len(body)))
		if err != nil {
			log.Fatal(err)
		}
	}
	err := os.Mkdir(destDirPath, 0755)
	if err != nil {
		log.Fatal(err)
	}

	// Read all the files from zip archive
	for _, f := range rp.File {

		// Store filename/path for returning and using later on
		fpath := filepath.Join(destDirPath, f.Name)

		fmt.Println("fpath", fpath) //for debug

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(destDirPath)+string(os.PathSeparator)) {
			return filesNames, fmt.Errorf("%s: illegal file path", fpath)
		}

		filesNames = append(filesNames, fpath)

		if f.FileInfo().IsDir() {
			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// Make File
		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return filesNames, err
		}

		rc, err := f.Open()
		if err != nil {
			return filesNames, err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return filesNames, err
		}

		_, err = io.Copy(outFile, rc)
		if err != nil {
			return filesNames, err
		}

		// Close the file without defer to close before next iteration of loop
		outFile.Close()
		rc.Close()

	}
	return filesNames, nil
}

// ExtractSome files from zip srcArchive , to dest directory
func ExtractSome(srcArchive string, isUrl bool, targetFiles []string, destDirPath string) error {

	var rp *zip.Reader

	if !isUrl { //local archive
		r_c, err := zip.OpenReader(srcArchive)
		if err != nil {
			log.Fatal(err)
		}
		defer r_c.Close()

		//create dest_dir_path for extractedFiles
		err = os.Mkdir(destDirPath, 0755)
		if err != nil {
			log.Fatal(err)
		}

		r := zip.Reader(r_c.Reader) //temp
		rp = &r

	} else if isUrl { //url archive

		// mk directory for unarchive files
		err := os.Mkdir(destDirPath, 0755)
		if err != nil {
			log.Fatal(err)
		}

		resp, err := http.Get(srcArchive)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		rp, err = zip.NewReader(bytes.NewReader(body), int64(len(body)))
		if err != nil {
			log.Fatal(err)
		}
	}

	for _, target := range targetFiles {

		for _, f := range rp.File {

			if f.Name != target {
				continue
			}

			rc, err := f.Open()
			if err != nil {
				log.Fatal(err)
			}

			fileName := filepath.Base(f.Name)
			out := filepath.Join(destDirPath, fileName)
			outFile, err := os.OpenFile(out, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())

			if err != nil {
				log.Fatal(err)
			}

			_, err = io.Copy(outFile, rc)
			if err != nil {
				log.Fatal(err)
			}

			outFile.Close()
			rc.Close()
			break
		}
	}
	return nil
}

// get metadata of zip srcArchive files. Set Fid to each directory/file.
func Meta(srcArchive string, isUrl bool) (zipMetadata []FileMetadata, err error) {

	var rp *zip.Reader
	filesInArch = 0
	// allocate dynamic array, start with minimum size
	FilesMetadata := make([]FileMetadata, 1)

	if !isUrl { //local archive

		r_c, err := zip.OpenReader(srcArchive)
		if err != nil {
			return nil, err
		}
		defer r_c.Close()

		r := zip.Reader(r_c.Reader) //temp
		rp = &r

	} else if isUrl { //url archive

		resp, err := http.Get(srcArchive)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		rp, err = zip.NewReader(bytes.NewReader(body), int64(len(body)))
		if err != nil {
			log.Fatal(err)
		}
	}

	for i, f := range rp.File {

		filesInArch++
		FilesMetadata[i].Fid = i
		FilesMetadata[i].FileName = f.Name
		FilesMetadata[i].CRC32 = f.CRC32
		FilesMetadata[i].UcSize = f.UncompressedSize64
		FilesMetadata[i].CSize = f.CompressedSize64
		FilesMetadata[i].CompressMethod = f.Method
		dataOffset, _ := f.DataOffset()

		FilesMetadata[i].StartInArch = uint64(dataOffset)
		if f.Method == zip.Deflate {
			FilesMetadata[i].EndInArch = uint64(dataOffset) + f.CompressedSize64
		} else {
			FilesMetadata[i].EndInArch = uint64(dataOffset) + f.UncompressedSize64
		}
		// resize dynamic array
		FilesMetadata = append(FilesMetadata, FilesMetadata[i])

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
