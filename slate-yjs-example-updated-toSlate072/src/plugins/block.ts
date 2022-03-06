import { Editor, Transforms, Element } from "slate";

const LIST_TYPES: string[] = ["numbered-list", "bulleted-list"];

export const toggleBlock = (editor: Editor, format: string): void => {
  const isActive = isBlockActive(editor, format);
  const isList = LIST_TYPES.includes(format);

  Transforms.unwrapNodes(editor, {
    match: (n) =>
      !Editor.isEditor(n) &&
      Element.isElement(n) &&
      LIST_TYPES.includes(n.type),
    split: true,
  });

  // Transforms.setNodes<Element>(editor, {
  //   type: isActive ? "paragraph" : isList ? "list-item" : format,
  // });
  const newProperties: Partial<Element> = {
    type: isActive ? "paragraph" : isList ? "list-item" : format,
  };
  Transforms.setNodes<Element>(editor, newProperties);

  if (!isActive && isList) {
    const block: any = { type: format, children: [] };
    Transforms.wrapNodes(editor, block);
  }
};

export const isBlockActive = (editor: Editor, format: string): boolean => {
  const [match] = Editor.nodes(editor, {
    match: (n) =>
      !Editor.isEditor(n) && Element.isElement(n) && n.type === format,
  });

  return !!match;
};
