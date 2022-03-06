/*eslint-disable */
import styled from "@emotion/styled";
import React, { useEffect, useMemo, useState } from "react";
import { createEditor, Node } from "slate";
import { withHistory } from "slate-history";
import { ReactEditor, withReact } from "slate-react";
import {
  SyncElement,
  toSharedType,
  useCursors,
  withCursor,
  withYjs,
} from "slate-yjs";
import { WebsocketProvider } from "y-websocket";
import * as Y from "yjs";
import { Button, H4, Instance, Title } from "./Components";

import { withLinks } from "./config/link";
import randomColor from "randomcolor";
import { Descendant } from "slate";
import type { CustomEditor } from "./types";
import { AnyObject, createPlateEditor, ImageToolbarButton, PlateEditor, TDescendant, TNode, withPlate } from "@udecode/plate";


// plate imports
import { ClientFrame} from "./Components";
import { CONFIG } from "./config/config";

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
  
} from "@udecode/plate";

import {
  AlignToolbarButtons,
  BasicElementToolbarButtons,
  BasicMarkToolbarButtons,
  IndentToolbarButtons,
  ListToolbarButtons,
  TableToolbarButtons,
} from "./config/components/Toolbars";

import { withStyledPlaceHolders } from "./config/components/withStyledPlaceHolders";
import {
  Check,
  FontDownload,
  FormatColorText,
  LineWeight,
  Image,
  Link,
  OndemandVideo,
} from "@styled-icons/material";


const WEBSOCKET_ENDPOINT =
  process.env.NODE_ENV === "production"
    ? "wss://demos.yjs.dev/slate-demo"
    : "ws://localhost:1234";

interface ClientProps {
  name: string;
  id: string;  
  slug: string;  
  removeUser: (id: any) => void;
}

const Client: React.FC<ClientProps> = ({ id, name, slug, removeUser }) => {

  const [value, setValue] = useState<TDescendant[]>([]);
  const [isOnline, setOnlineState] = useState<boolean>(false);

  const onChange= (value:TDescendant[]) => setValue(value);

  const color = useMemo(
    () =>
      randomColor({
        luminosity: "dark",
        format: "rgba",
        alpha: 1,
      }),
    []
  );

  const [sharedType, provider] = useMemo(() => {
    const doc = new Y.Doc();
    const sharedType = doc.getArray<SyncElement>("content");
    const provider = new WebsocketProvider(WEBSOCKET_ENDPOINT, slug, doc, {
      connect: false,
    });
    console.log('doc=',doc); // for debug
    return [sharedType, provider];
  }, [id]);
  
  const editableProps = {
    placeholder: "Type…",
        
  }; 

  let components = createPlateUI();
  components = withStyledPlaceHolders(components);

  const plugins = createPlugins(
    [
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
      
    ],
    {
      components,
    }
  );


  const editor = useMemo(() => {
    const editor = withCursor(
      withYjs(
        withLinks(withReact(withHistory(createPlateEditor({id:id , plugins})))), 
        sharedType
      ),
      provider.awareness
    );
    

    return editor;
  }, [sharedType, provider]);

  useEffect(() => {
    provider.on("status", ({ status }: { status: string }) => {
      setOnlineState(status === "connected");
    });

    provider.awareness.setLocalState({
      alphaColor: color.slice(0, -2) + "0.2)",
      color,
      name,
    });

    // Provide a initial value from the client.
    // In prod, when using y-websocket - provide initial state from the server.
    
    provider.on("sync", (isSynced: boolean) => {
      if (isSynced && sharedType.length === 0) {
        toSharedType(sharedType, [
          { type: "paragraph", children: [{ text: "xxx" }] },
        ]);
      }
    });

    provider.connect();

    return () => {
      provider.disconnect();
    };
  }, [provider]);

  const { decorate } = useCursors(editor);

  const toggleOnline = () => {
    isOnline ? provider.disconnect() : provider.connect();
  };

  return (
    <Instance online={isOnline}>
      <Title>
        <Head>Editor: {name}</Head>
        <div style={{ display: "flex", marginTop: 10, marginBottom: 10 }}>
          <Button type="button" onClick={toggleOnline}>
            Go {isOnline ? "offline" : "online"}
          </Button>
          <Button type="button" onClick={() => removeUser(id)}>
            Remove
          </Button>
        </div>
      </Title>

      <ClientFrame>
      
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
            tooltip={{ content: "Text color" }}
          />
          <ColorPickerToolbarDropdown
            pluginKey={MARK_BG_COLOR}
            icon={<FontDownload />}
            selectedIcon={<Check />}
            tooltip={{ content: "Highlight color" }}
          />
          <AlignToolbarButtons />
          <LineHeightToolbarDropdown icon={<LineWeight />} />
          <LinkToolbarButton icon={<Link />} />
          <ImageToolbarButton icon={<Image />} />
          <MediaEmbedToolbarButton icon={<OndemandVideo />} />
          <TableToolbarButtons />
        </HeadingToolbar>     

        <Plate
          id={id}
          editableProps={editableProps} 
          //initialValue={initialValue}
          //plugins={plugins}

          editor={editor}
          value={value}
          onChange={onChange}
          decorate={decorate}
          
          
          
        ></Plate>
      </>
    </ClientFrame>
    </Instance>
  );
};

export default Client;

const Head = styled(H4)`
  margin-right: auto;
`;
