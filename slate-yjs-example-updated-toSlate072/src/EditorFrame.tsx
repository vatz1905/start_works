import Caret from "./Caret";
import React, { useCallback } from "react";
import { Node } from "slate";
import {
  Editable,
  ReactEditor,
  RenderLeafProps,
  Slate,
  useSlate,
} from "slate-react";
import { ClientFrame, Icon, IconButton } from "./Components";
import { isBlockActive, toggleBlock } from "./plugins/block";
import { insertLink, isLinkActive, unwrapLink } from "./plugins/link";
import { isMarkActive, toggleMark } from "./plugins/mark";
import { Descendant } from "slate";

export interface EditorFrame {
  editor: ReactEditor;
  value: Descendant[];
  onChange: (value: Descendant[]) => void;
  decorate: any;
}

const renderElement = (props: any) => <Element {...props} />;

const EditorFrame: React.FC<EditorFrame> = ({
  editor,
  value,
  onChange,
  decorate,
}) => {
  const renderLeaf = useCallback(
    (props: any) => <Leaf {...props} />,

    [decorate]
  );
  console.log("editor=", editor); //for debug
  return (
    <ClientFrame>
      <Slate editor={editor} value={value} onChange={onChange}>
        <div
          style={{
            display: "flex",
            flexWrap: "wrap",
            position: "sticky",
            top: 0,
            backgroundColor: "white",
            zIndex: 1,
          }}
        >
          <MarkButton format="bold" icon="format_bold" />
          <MarkButton format="italic" icon="format_italic" />
          <MarkButton format="underline" icon="format_underlined" />
          <MarkButton format="code" icon="code" />

          <BlockButton format="heading-one" icon="looks_one" />
          <BlockButton format="heading-two" icon="looks_two" />
          <BlockButton format="block-quote" icon="format_quote" />

          <BlockButton format="numbered-list" icon="format_list_numbered" />
          <BlockButton format="bulleted-list" icon="format_list_bulleted" />

          <LinkButton />
        </div>

        <Editable
          renderElement={renderElement}
          renderLeaf={renderLeaf}
          decorate={decorate}
        />
      </Slate>
    </ClientFrame>
  );
};

export default EditorFrame;

const Element: React.FC<any> = ({ attributes, children, element }) => {
  switch (element.type) {
    case "link":
      return (
        <a {...attributes} href={element.href}>
          {children}
        </a>
      );
    case "block-quote":
      return <blockquote {...attributes}>{children}</blockquote>;
    case "bulleted-list":
      return <ul {...attributes}>{children}</ul>;
    case "heading-one":
      return <h1 {...attributes}>{children}</h1>;
    case "heading-two":
      return <h2 {...attributes}>{children}</h2>;
    case "list-item":
      return <li {...attributes}>{children}</li>;
    case "numbered-list":
      return <ol {...attributes}>{children}</ol>;
    default:
      return <p {...attributes}>{children}</p>;
  }
};

const Leaf: React.FC<RenderLeafProps> = ({ attributes, children, leaf }) => {
  if (leaf.bold) {
    children = <strong>{children}</strong>;
  }

  if (leaf.code) {
    children = <code>{children}</code>;
  }

  if (leaf.italic) {
    children = <em>{children}</em>;
  }

  if (leaf.underline) {
    children = <u>{children}</u>;
  }

  const data = leaf.data as any;

  return (
    <span
      {...attributes}
      style={
        {
          position: "relative",
          backgroundColor: data?.alphaColor,
        } as any
      }
    >
      {leaf.isCaret ? <Caret {...(leaf as any)} /> : null}
      {children}
    </span>
  );
};

const BlockButton: React.FC<any> = ({ format, icon }) => {
  const editor = useSlate();
  return (
    <IconButton
      active={isBlockActive(editor, format)}
      onMouseDown={(event: React.MouseEvent) => {
        event.preventDefault();
        toggleBlock(editor, format);
      }}
    >
      <Icon className="material-icons">{icon}</Icon>
    </IconButton>
  );
};

const MarkButton: React.FC<any> = ({ format, icon }) => {
  const editor = useSlate();
  return (
    <IconButton
      active={isMarkActive(editor, format)}
      onMouseDown={(event: React.MouseEvent) => {
        event.preventDefault();
        toggleMark(editor, format);
      }}
    >
      <Icon className="material-icons">{icon}</Icon>
    </IconButton>
  );
};

const LinkButton = () => {
  const editor = useSlate();

  const isActive = isLinkActive(editor);

  return (
    <IconButton
      active={isActive}
      onMouseDown={(event: React.MouseEvent) => {
        event.preventDefault();

        if (isActive) return unwrapLink(editor);

        const url = window.prompt("Enter the URL of the link:");

        url && insertLink(editor, url);
      }}
    >
      <Icon className="material-icons">link</Icon>
    </IconButton>
  );
};
