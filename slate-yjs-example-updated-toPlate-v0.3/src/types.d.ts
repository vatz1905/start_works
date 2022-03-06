import { CursorEditor, YHistoryEditor, YjsEditor } from '@slate-yjs/core';
import { BaseEditor } from 'slate';
import { Descendant } from 'slate';
import { HistoryEditor } from 'slate-history';
import { ReactEditor } from 'slate-react';

export type CursorData = {
  name: string;
  color: string;
};

// export type EmptyText = {
//   text: string
// }
export type Type = {
  type: string
}

export type CustomText = {
  bold?: boolean;
  italic?: boolean;
  underline?: boolean;
  strikethrough?: boolean;
  code?:boolean;
  data?:boolean;
  text: string;
  isCaret?:boolean;
};

export type Paragraph = {
  type: 'paragraph';
  children: Descendant[];
};

export type InlineCode = {
  type: 'inline-code';
  children: Descendant[];
};

export type HeadingOne = {
  type: 'heading-one';
  children: Descendant[];
};

export type HeadingTwo = {
  type: 'heading-two';
  children: Descendant[];
};

export type BlockQuote = {
  type: 'block-quote';
  children: Descendant[];
};

export type BulletedList = {
  type: 'bulleted-list';
  children: Descendant[];
};

export type NumberedList = {
  type: 'numbered-list';
  children: Descendant[];
};

export type ListItem = {
  type: 'list-item';
  children: Descendant[];
};

export type Link = {
  type: 'link';
  url: string;
  children: Descendant[];
};

export type CustomElement =
  | Paragraph
  | InlineCode
  | HeadingOne
  | HeadingTwo
  | BlockQuote
  | BulletedList
  | NumberedList
  | ListItem
  | Link
  | Type;

export type CustomEditor = BaseEditor & ReactEditor & HistoryEditor; 

declare module 'slate' {
  interface CustomTypes {
    Editor: CustomEditor;
    Element: CustomElement;
    Text: CustomText;
  }
}

