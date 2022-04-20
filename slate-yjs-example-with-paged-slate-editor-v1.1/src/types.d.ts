import { ReactEditor } from 'slate-react';
import { HistoryEditor } from 'slate-history';
import { BaseEditor, Descendant } from 'slate';
import { CursorEditor, YHistoryEditor, YjsEditor } from '@slate-yjs/core';
import { string } from 'lib0';

export type CustomText = {  
  text: string;  
};

export type PageElement = {
  type: 'page';
  children: Descendant[];
};
export type Paragraph = {
  type: 'paragraph';
  children: Descendant[];
};

export type CustomElement = PageElement | Paragraph;  

export type CustomEditor = ReactEditor &
  BaseEditor &
  HistoryEditor &
  YjsEditor &
  YHistoryEditor &
  CursorEditor<CursorData>;

declare module 'slate' {
  interface CustomTypes {
    Editor: CustomEditor;
    Element: CustomElement;
    Text: CustomText;
  }
}
