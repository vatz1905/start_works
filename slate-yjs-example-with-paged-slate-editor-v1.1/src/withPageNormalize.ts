/*eslint-disable */

import { ReactEditor } from 'slate-react';
import { Element, Text, Transforms,} from 'slate';
import { CustomEditor } from './types';

//
let pages: number;
let pageEditor: HTMLElement;
let pageNode: HTMLElement; 
let pageHeight: number;
let currEditPage: number;
let lastPage: number;
let childLength: number;
let currPageOverflow: boolean = false;
// 
let lineDeleted: boolean = false;  //temp true for testing     
let childLengthCurrentValue: number = 0;
let childLengthBeforeLineDeleted: number = 0; // changed
//let testFlag=false;  //temp flag for testing 
let pHeight : number;   // temp for testing
let pElementsHeightArr : number[]; // temp for testing

export const withPageNormalize = ( editor:ReactEditor ) => { 

  const { normalizeNode } = editor; 

  editor.normalizeNode = (entry) => {      


      const [ node, path ] = entry;

      console.log( "in entry" );  // for debug 


      if ( Text.isText( node ) ) return normalizeNode( entry ); 

      try {
          pageEditor= ReactEditor.toDOMNode( editor, editor ); 
          //console.log( 'pageEditor=', pageEditor ); //for debug
                              
      } catch ( err ) {

          console.log( "pageNode err" );  // for debug 
          return;
          //return normalizeNode(entry);
      }



      
      if ( Element.isElement( node ) && node.type === 'page') {  // todo -   && node.type === 'page' ?

          console.log( "in if is element" );  // for debug 
         
          //console.log('editor.children=', editor.children );
         

          try {

              pageNode= ReactEditor.toDOMNode( editor, node );  
              //console.log( 'pageNode=', pageNode ); // for debug               
             
          } catch ( err ) {                
              return;
              //return normalizeNode(entry);
          }            

          pages=editor.children.length;
          currEditPage = path[0];
          lastPage = pages-1;

          

        //   console.log( 'pages=', pages ); // for debug
        //   console.log( 'lastPage=', lastPage ); // for debug
        //   console.log( 'currEditPage=',currEditPage ); //for debug
          

        //   console.log( 'path=', path ); //for debug 

          pageHeight = getPageDeclaredHeight( pageNode );
        //   console.log( 'pageHeight=', pageHeight ); // for debug    
          
          
          childLength = pageNode.children.length;  
        //   console.log( 'childLength=', childLength ); //for debug

          //for checking if delete state
          childLengthCurrentValue = node.children.length;
         
          console.log( 'childLengthCurrentValue=', childLengthCurrentValue ); //for debug

          console.log( 'childLengthBeforeLineDeleted=', childLengthBeforeLineDeleted ); //for debug

          //check if line was deleted
          if ( childLengthBeforeLineDeleted > childLengthCurrentValue && currEditPage !== lastPage ) { //changed to >=

              lineDeleted = true;

              console.log( 'lineDeleted=', lineDeleted ); //for debug
          }
          
          
        

          //edit the last new page

          if ( currEditPage === lastPage ) { 

              console.log('in edit new page'); //for debug 

              // check if current last page is overflowed
              currPageOverflow = isPageOverflow( pageNode, pageHeight );

              [ pHeight, pElementsHeightArr ]=pageElementsHeight(pageNode); // temp for testing
            //   console.log('pHeight=',pHeight); //for debug  // temp for testing
            //   console.log('pElementsHeightArr=',pElementsHeightArr); //for debug  // temp for testing

            //   console.log('openPage=',currPageOverflow); //for debug 

              if ( currPageOverflow ) { 

                  Transforms.splitNodes( editor, { at:[ path[0], childLength-1 ]}); 

                  lastPage = lastPage+1;

                  currEditPage = currEditPage+1; // to prevent delayed update of curr_edit_page (to the next if)
              }
                 
          }

          


          //edit exist page (except last page) - add new lines/elements and edit lines/elements

          if ( !lineDeleted && currEditPage !== lastPage ) {

             console.log( 'in edit exist page' ); //for debug

          //    let overflow = false; //init
             
          //     overflow = isPageOverflow( pageNode, pageHeight );  
              
          //     if ( overflow ) { 

          //         Transforms.moveNodes( editor, {          
          //             at: [ path[0], childLength ],    
          //             to: [ path[0]+1, 0 ], 
          //         });

                  //overflow = false; // init before iteration 

                 
                  
                  for ( let i = currEditPage; i <= lastPage; i++ ) { 

                      //local vars and init
                      let nextPageHeight = 0;     
                      let nextPageChildLength = 0; 
                      let nextPageNode = pageEditor.children[i] as HTMLElement;
                      let nextOverflow = false;
                      
                      // console.log( 'nextPageNode=', nextPageNode ); //for debug 

                      nextPageChildLength = nextPageNode.children.length;
                      //console.log( 'nextPageChildLength=', nextPageChildLength ); //for debug 

                      nextPageHeight = getPageDeclaredHeight( nextPageNode ); 
                      //console.log( 'nextPageHeight=', nextPageHeight ); //for debug 
                      
                      nextOverflow = isPageOverflow( nextPageNode, nextPageHeight );
                      console.log( 'overflow=', nextOverflow ); //for debug 

                      if ( !nextOverflow ) {

                          console.log( 'in tp- break' ); //for debug
                         
                          break;
                      }

                      if ( nextOverflow && i !== lastPage ) { 
                          
                          Transforms.moveNodes( editor, {          
                              at: [ i, nextPageChildLength ],  
                              to: [ i+1, 0 ], 
                          });
                          // init overflow
                          nextOverflow = false; 
                      }
                      else if ( nextOverflow && i === lastPage ) { 
                          
                          Transforms.splitNodes( editor, { at:[ i, nextPageChildLength-1 ]});
                            
                          // Transforms.moveNodes( editor, {          
                          //     at: [ i, nextPageChildLength ],   
                          //     to: [ i+1, 0 ], 
                          // });
                          //lastPage = lastPage+1; 
                          
                          // init overflow
                          nextOverflow = false; 
                         
                      }
                  }
              //}  
          }

          // edit exist page (except last page) - delete line

          if ( lineDeleted && currEditPage !== lastPage ) {

              console.log( 'in edit deleted line' ); //for debug

              //init deleted_line flag
              lineDeleted = false;     

              //local vars
              let firstNextPageNode : HTMLElement;
              let firstNextPageChildLength : number;
              let firstNextPageHeight : number;
              let firstNextPageElementsHeightArr : number[];
              let secondNextPageNode : HTMLElement;
              let secondNextPageChildLength : number;  
              let secondNextPageHeight : number;
              let secondNextPageElementsHeightArr : number[];
              // temp
              let lastPageNode : HTMLElement;
              let lastPageChildLength : number;
              let lastPageHeight : number;
              let lastPageElementsHeightArr : number[];
              
              

              let elementsHeight : number;
              let lastIndex : number;
              let elementsHeightToMove : number = 0;
              let maxElementsHeight : number; 

              

              for ( let i = currEditPage; i <= lastPage-1; i++ ) {



                  //init vars before iteration
                  firstNextPageChildLength = 0;  
                  firstNextPageHeight = 0;
                  firstNextPageElementsHeightArr = [];                    
                  secondNextPageChildLength = 0;  
                  secondNextPageHeight = 0;
                  secondNextPageElementsHeightArr = [];
                  lastPageChildLength = 0;
                  lastPageHeight = 0;
                  lastPageElementsHeightArr = [];
                  

                  //
                  firstNextPageNode = pageEditor.children[i] as HTMLElement;

                  firstNextPageChildLength = firstNextPageNode.children.length;

                  [ firstNextPageHeight, firstNextPageElementsHeightArr ] = pageElementsHeight( firstNextPageNode );

                  //maxElementsHeight to move from next (second) page to this (first) page
                  maxElementsHeight = pageHeight - firstNextPageHeight// + elementsHeightToMove;

                  //init
                  elementsHeight = 0;
                  lastIndex = 0;

                  if ( i !== lastPage-1 ) {

                      secondNextPageNode = pageEditor.children[i+1] as HTMLElement;

                      secondNextPageChildLength = firstNextPageNode.children.length;

                      [ secondNextPageHeight, secondNextPageElementsHeightArr ] = pageElementsHeight( secondNextPageNode );
                  

                      // check how many elements to move from next (second) page to current page
                      // if last_index=-1 no element have to be moved.                        
                      for ( let j = 0; j < secondNextPageChildLength; j++ ) {

                          elementsHeight = elementsHeight + secondNextPageElementsHeightArr[j];

                          if ( elementsHeight > maxElementsHeight ) {
                              lastIndex = j-1;
                              break;
                          }
                          //elementsHeightToMove = elementsHeight;
                      }
                      if ( lastIndex !== -1 ) {

                          // moveNodes from page i+1 to page i
                          for ( let k = 0; k <= lastIndex; k++ ) {

                              if ( i !== lastPage-1 ) {
                                  Transforms.moveNodes( editor, {          
                                      at: [ i+1, 0 ],  
                                      to: [ i, firstNextPageChildLength+k ],  // todo
                                  });
                                  
                              }
                          }
                      }    
                  }
                  else if ( i === lastPage-1 ) {

                      lastPageNode = pageEditor.children[i+1] as HTMLElement;
                      
                      lastPageChildLength = lastPageNode.children.length;

                      [ lastPageHeight, lastPageElementsHeightArr ] = pageElementsHeight( lastPageNode );

                      // check how many elements to move from next (last) page to last-1 page
                      // if last_index=-1 no element have to be moved.                        
                      for ( let j = 0; j < lastPageChildLength; j++ ) {

                          elementsHeight = elementsHeight + lastPageElementsHeightArr[j];

                          if ( elementsHeight > maxElementsHeight ) {
                              lastIndex = j-1;
                              break;
                          }
                          //elementsHeightToMove = elementsHeight;
                      }
                      if ( lastIndex !== -1 ) {

                          // moveNodes from page i+1 to page i
                          for ( let k = 0; k <= lastIndex; k++ ) {
                              
                              
                              if ( lastPageChildLength > 1 ) {
                                          
                                  Transforms.moveNodes( editor, {          
                                      at: [ i+1, 0 ],  
                                      to: [ i, firstNextPageChildLength+k ],  // todo
                                  });
                                          
                              }
                              else if ( lastPageChildLength === 1 ) {

                                  Transforms.moveNodes( editor, {          
                                      at: [ i+1, 0 ],  
                                      to: [ i, firstNextPageChildLength+k ],  //  todo 
                                  });

                                  Transforms.removeNodes( editor, { at:[ i+1 ]});
                                          
                                  lastPage = lastPage-1;
                                  lineDeleted = false; //init      
                                  break; // only one element left
                                  }
                              }
                          }
                  }
                  else {
                      lineDeleted = false; //init  
                      break;

                  }   

              } 

          }
       
      }
     

      // save last child_len value for the next iteration to check if line was deleted
      childLengthBeforeLineDeleted=childLengthCurrentValue ;

      //console.log( 'childLengthBeforeLineDeleted=', childLengthBeforeLineDeleted )
               
      return normalizeNode( entry );       
  }; 

  return editor;
};





// helper functions

function getPageDeclaredHeight( page:HTMLElement ) {

  const style = window.getComputedStyle( page );
  const computedHeight = page.offsetHeight;
  const padding = parseFloat( style.paddingLeft ) + parseFloat( style.paddingRight );
  let pageHeight = computedHeight - padding;
  //console.log( 'pageHeight=', pageHeight ) // for debug
  return pageHeight;
}

//
function isPageOverflow( page : HTMLElement, pageHeight : number ) : boolean { 

    console.log( 'in isPageOverflow func' ); //for debug
    let pageOverflow = false;
    let currPageHeight = 0;

    const children = Array.from( page.children );        
            
    children.forEach( ( child, childPath )  => {
                                  
      const childStyles = window.getComputedStyle( child );
      const computedChildHeight = child.clientHeight;
      const childMargin =parseFloat( childStyles.marginTop ) //+ parseFloat( childStyles.marginBottom ); //temp change
      const childPadding = parseFloat( childStyles.paddingTop ) + parseFloat( childStyles.paddingBottom );
      const childBorder = parseFloat( childStyles.borderLeftWidth ) + parseFloat( childStyles.borderRightWidth )+ parseFloat( childStyles.borderTopWidth ) + parseFloat( childStyles.borderBottomWidth );
      //const childBorder = parseFloat( childStyles.borderTopWidth ) + parseFloat( childStyles.borderBottomWidth );                                
      const childHeight = computedChildHeight + childMargin + childPadding + childBorder;
      //console.log( 'computedChildHeight=', computedChildHeight ); //for debug
      console.log( 'childHeight=', childHeight ); //for debug
      //console.log( 'childMargin=', childMargin ); //for debug                                
      currPageHeight = currPageHeight + childHeight;
      console.log( 'currentpageHeight=', currPageHeight ); //for debug 
    
      // check if current page is overflowed 
       
      if ( currPageHeight > pageHeight ) { 

          return pageOverflow = true;      
                  
      } 

  });

  return pageOverflow;
}
//
function pageElementsHeight( page : HTMLElement ) : ( any[] ) { 

    console.log( 'in pageElementsHeight' ); //for debug
    
    let accumulatePageHeight = 0;
    let elementIndex = 0;
    let elementsHeight : number[] = [];

    const children = Array.from( page.children );        
            
    children.forEach( ( child, childPath )  => {
                                  
      const childStyles = window.getComputedStyle( child );
      const computedChildHeight = child.clientHeight;
      const childMargin = parseFloat( childStyles.marginTop )//+ parseFloat( childStyles.marginBottom ); //temp change
      const childPadding = parseFloat( childStyles.paddingTop ) + parseFloat( childStyles.paddingBottom );
      const childBorder = parseFloat( childStyles.borderLeftWidth ) + parseFloat( childStyles.borderRightWidth )+ parseFloat( childStyles.borderTopWidth ) + parseFloat( childStyles.borderBottomWidth );
      //const childBorder = parseFloat( childStyles.borderTopWidth ) + parseFloat( childStyles.borderBottomWidth );
                                      
      const childHeight = computedChildHeight + childMargin + childPadding + childBorder;
      //
      elementsHeight[ elementIndex++ ] = childHeight;
      //console.log( 'elementsHeight[]=', elementsHeight ); //for debug
      //console.log( 'childMargin=', childMargin ); //for debug 
      accumulatePageHeight = accumulatePageHeight + childHeight;
      console.log( 'currentpageHeight=', accumulatePageHeight ); //for debug 
      
      
  });
  //console.log( 'elementsHeight[]=', elementsHeight ); //for debug

  return   [ accumulatePageHeight , elementsHeight ];
}
