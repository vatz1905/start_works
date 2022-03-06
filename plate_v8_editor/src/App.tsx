import React, {  useState} from 'react';
import { CONFIG } from './config/config';
//import { PLUGINS } from './config/plugins';
import {   
  createParagraphPlugin,
  createBlockquotePlugin,
  createTodoListPlugin,
  createHeadingPlugin,
  createImagePlugin,
  createHorizontalRulePlugin,
  createLineHeightPlugin,
  createLinkPlugin,
  createListPlugin,
  createTablePlugin,
  createMediaEmbedPlugin,
  //createExcalidrawPlugin,
  createCodeBlockPlugin,
  createAlignPlugin,
  createBoldPlugin,
  createCodePlugin,
  createItalicPlugin,
  createHighlightPlugin,
  createUnderlinePlugin,
  createStrikethroughPlugin,
  createSubscriptPlugin,
  createSuperscriptPlugin,
  createFontBackgroundColorPlugin,
  createFontFamilyPlugin,
  createFontColorPlugin,
  createFontSizePlugin,
  createFontWeightPlugin,
  createKbdPlugin,
  createNodeIdPlugin,
  createIndentPlugin,
  createAutoformatPlugin,
  createResetNodePlugin,
  createSoftBreakPlugin,
  createExitBreakPlugin,
  createNormalizeTypesPlugin,
  createTrailingBlockPlugin,
  createSelectOnBackspacePlugin,
  createComboboxPlugin,
  createMentionPlugin,
  createDeserializeMdPlugin,
  createDeserializeCsvPlugin,
  createDeserializeDocxPlugin,
  createJuicePlugin,
  createPlateUI,
  createPlugins,
  HeadingToolbar,
  Plate,
  ColorPickerToolbarDropdown,
  MARK_COLOR,
  MARK_BG_COLOR,
  LineHeightToolbarDropdown,
  LinkToolbarButton,
  //ImageToolbarButton,
  MediaEmbedToolbarButton,
  //MentionCombobox,

} from '@udecode/plate'



import {AlignToolbarButtons, BasicElementToolbarButtons, BasicMarkToolbarButtons, IndentToolbarButtons, ListToolbarButtons, TableToolbarButtons} from './config/components/Toolbars';
import { withStyledPlaceHolders } from './config/components/withStyledPlaceHolders';
import { Check, FontDownload, FormatColorText, LineWeight, Link, OndemandVideo } from '@styled-icons/material';
//import { VALUES } from './config/values/values';


//let node=createNode('page')
//const styledPage=StyledElement( )
//node= {type:'page',children:[{type:'paragraph',children:[{type:'text',text:''}]}]}

const pageNode= {
  //type:node.type,
  //children: [Element],
  style: {
     
    padding: '60px',
    marginLeft: 'auto',
    marginRight: 'auto',
    marginTop: '40px',
    marginBottom: '40px',
    width: '816px',
    height: '1056px',    
    background:'#eee',
  
  },
}



function App() {
   
  const editableProps = {   
    placeholder: 'Type…',
    style: pageNode.style,   
  };
  
  
   //temp
//   const emptyPage ={type:'page',children:[{type:'paragraph',children:[{type:'text',text:''}]}]}
   //end temp

//   const pageNode=createDocumentNode()
  
   
   const initialValue = [  
  
    {
      type:'page',//pageNode.type,
      children: [
        {
          type: 'paragraph',
          children: [
            {
              type: 'text',
              text: 'Hey there! '
            },        
          ]
        },
        {
          type: 'paragraph',
          children: [
            {
              type: 'text',
              text: 'This is first page of whatever you want to write! continue writing!'
            },
          //  {
          //    type: 'text',
          //    text: 'anything you want and download the button on top!'
           // },
          ]
        },
      ]
    }
  ]    
  
  
  
  
  
  
  // const editor:any=usePlateEditorRef();
  // const currNode=getNodes(editor) //return generator - todo 
  // console.log('currNode=',currNode)
  




  const [debugValue, setDebugValue] = useState(null);
  
  
 
  
  
  let components = createPlateUI();
  components = withStyledPlaceHolders(components);


  const plugins = createPlugins([
    createParagraphPlugin(),
    createBlockquotePlugin(),
    createTodoListPlugin(),
    createHeadingPlugin(),
    createImagePlugin(),
    createHorizontalRulePlugin(),
    createLineHeightPlugin(),
    createLinkPlugin(),
    createListPlugin(),
    createTablePlugin(),
    createMediaEmbedPlugin(),    
    createCodeBlockPlugin(),
    createAlignPlugin(),
    createBoldPlugin(),
    createCodePlugin(),
    createItalicPlugin(),
    createHighlightPlugin(),
    createUnderlinePlugin(),
    createStrikethroughPlugin(),
    createSubscriptPlugin(),
    createSuperscriptPlugin(),
    createFontBackgroundColorPlugin(),
    createFontFamilyPlugin(),
    createFontColorPlugin(),
    createFontSizePlugin(),
    createFontWeightPlugin(),
    createKbdPlugin(),
    createNodeIdPlugin(),
    createIndentPlugin(),
    createAutoformatPlugin(),
    createResetNodePlugin(CONFIG.resetBlockType),
    createSoftBreakPlugin(CONFIG.softBreak),
    createExitBreakPlugin(CONFIG.exitBreak),
    createNormalizeTypesPlugin(),
    createTrailingBlockPlugin(CONFIG.trailingBlock),
    createSelectOnBackspacePlugin(CONFIG.selectOnBackspace),
    createComboboxPlugin(),
    createMentionPlugin(),
    createDeserializeMdPlugin(),
    createDeserializeCsvPlugin(),
    createDeserializeDocxPlugin(),
    createJuicePlugin(),

  
  ], {
    components,
  });
  
  

  return (
    <>
      <HeadingToolbar>
        <BasicElementToolbarButtons />
        <ListToolbarButtons />
        <IndentToolbarButtons />
        <BasicMarkToolbarButtons />
        <ColorPickerToolbarDropdown
          pluginKey={MARK_COLOR}
          icon={<FormatColorText />}
          selectedIcon={<Check />}
          tooltip={{ content: 'Text color' }}
        />
        <ColorPickerToolbarDropdown
          pluginKey={MARK_BG_COLOR}
          icon={<FontDownload />}
          selectedIcon={<Check />}
          tooltip={{ content: 'Highlight color' }}
        />
        <AlignToolbarButtons />
        <LineHeightToolbarDropdown icon={<LineWeight />} />
        <LinkToolbarButton icon={<Link />} />
        
        <MediaEmbedToolbarButton icon={<OndemandVideo />} />
        <TableToolbarButtons />
      </HeadingToolbar>

      
      
      <Plate
        editableProps={editableProps } 
        initialValue={ initialValue}
        plugins={plugins}        
        
        onChange={(newValue:any) => {
          setDebugValue(newValue);
          // save newValue...
        }}
       
      >
      value: {JSON.stringify(debugValue)}
      </Plate>
    </>
  );
}


export default App;


//  <ImageToolbarButton icon={<Image />} />
