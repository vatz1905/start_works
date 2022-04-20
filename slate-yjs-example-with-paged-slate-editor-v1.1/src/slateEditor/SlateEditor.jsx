/*eslint-disable */
import React, { useCallback, useMemo } from 'react';
import {
  Editable,
  ReactEditor,
  RenderLeafProps,
  Slate,  
} from "slate-react";
import { Descendant } from "slate";
import './Editor.css';
import Caret from "../Caret";
import { ClientFrame } from "../Components";
import Toolbar from './Toolbar/Toolbar';
import { sizeMap, fontFamilyMap } from './utils/SlateUtilityFunctions.js';
import withTables from './plugins/withTable.js';
import withEmbeds from './plugins/withEmbeds.js';
import Link from'./Elements/Link/Link';
import Image from './Elements/Image/Image';
import Video from './Elements/Video/Video';
import { withPageNormalize } from '../withPageNormalize'; // change for pages


// export interface SlateEditor {
//   editor: ReactEditor;
//   value: Descendant[];
//   onChange: (value: Descendant[]) => void;
//   decorate: any;
// }

export const SlateEditor= ({  
  editor,
  value,
  onChange,
  decorate,
}) => { 

  editor = useMemo(() => withPageNormalize(withEmbeds(withTables(editor))), []);
  
  const renderElement = useCallback(props => <Element {...props}/>,[]);

  const renderLeaf = useCallback(
    (props) => <Leaf {...props} />,

    [decorate]
  );
 
  return (
    <ClientFrame >
      <Slate editor={editor} value={value} onChange={onChange}>
        <Toolbar />
        <div className="editor-wrapper" style={{border:'1px solid #f3f3f3' , padding:'0 10px'}}>
        
          <Editable
            placeholder='Write something'
            renderElement={renderElement}
            renderLeaf={renderLeaf}
            decorate={decorate}
          />
        </div>
      </Slate>
    </ClientFrame>
  );
};

export default SlateEditor;

const Element = (props) => {

  const {attributes, children, element} = props;

  switch (element.type) {
    case 'headingOne':
      return <h1 {...attributes}>{children}</h1>
    case 'headingTwo':
      return <h2 {...attributes}>{children}</h2>
    case 'headingThree':
      return <h3 {...attributes}>{children}</h3>
    case 'blockquote':
      return <blockquote {...attributes}>{children}</blockquote>
    case 'alignLeft':
      return <div style={{textAlign:'left',listStylePosition:'inside'}} {...attributes}>{children}</div>
    case 'alignCenter':
      return <div style={{textAlign:'center',listStylePosition:'inside'}} {...attributes}>{children}</div>
    case 'alignRight':
      return <div style={{textAlign:'right',listStylePosition:'inside'}} {...attributes}>{children}</div>
    case 'list-item':
      return  <li {...attributes}>{children}</li>
    case 'orderedList':
      return <ol type='1' {...attributes}>{children}</ol>
    case 'unorderedList':
      return <ul {...attributes}>{children}</ul>
    case 'link':
      return <Link {...props}/>
       
    case 'table':
      return  <table>
                <tbody {...attributes}>{children}</tbody>
              </table>
    case 'table-row':
      return <tr {...attributes}>{children}</tr>
    case 'table-cell':
      return <td {...attributes}>{children}</td>
    case 'image':
      return <Image {...props}/>
    case 'video':
      return <Video {...props}/>
    case 'page':                    // change for pages
      return <div className="page" {...attributes}  
        style={{
          position: 'relative',
          marginLeft: 'auto',
          marginRight: 'auto',
          marginTop: '40px',
          marginBottom: '40px',
          width: '816px',
          height: '1056px',
          padding: '60px',
          background:'#eee',
  
       }}>
          {children}
       </div>

    default :
      return <p {...attributes}>{children}</p>
  }
};

const Leaf = ({ attributes, children, leaf }) => {
  if (leaf.bold) {
    children = <strong>{children}</strong>
  }

  if (leaf.code) {
    children = <code>{children}</code>
  }

  if (leaf.italic) {
    children = <em>{children}</em>
  }
  if(leaf.strikethrough){
      children = <span style={{textDecoration:'line-through'}}>{children}</span>
  }
  if (leaf.underline) {
    children = <u>{children}</u>
  }
  if(leaf.superscript){
      children = <sup>{children}</sup>
  }
  if(leaf.subscript){
      children = <sub>{children}</sub>
  }
  if(leaf.color){
      children = <span style={{color:leaf.color}}>{children}</span>
  }
  if(leaf.bgColor){
      children = <span style={{backgroundColor:leaf.bgColor}}>{children}</span>
  }
  if(leaf.fontSize ){
      const size = (sizeMap )[leaf.fontSize ]
      children = <span style={{fontSize:size}}>{children}</span>
  }
  if(leaf.fontFamily ){
      const family = (fontFamilyMap )[leaf.fontFamily]
      children = <span style={{fontFamily:family}}>{children}</span>
  }
  

  const data = leaf.data ;

  return (
    <span
      {...attributes}
      style={
        {
          position: "relative",
          backgroundColor: data?.alphaColor,
        } 
      }
    >
      {leaf.isCaret ? <Caret {...(leaf )} /> : null}
      {children}
    </span>
  );
  
};

