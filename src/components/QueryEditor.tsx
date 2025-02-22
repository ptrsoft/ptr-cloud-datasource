import React, { useEffect, useRef, useState, useCallback } from 'react';
import { InlineField, Select, Input } from '@grafana/ui';
import {
  RESPONSE_TYPE,
  getCloudElementsQuery
} from '../common-ds';
import { EditorRow, EditorRows } from '../extended/EditorRow';
import { Services } from '../service';

export function QueryEditor({ query, onChange, onRunQuery, datasource }: any) {
  const defaultFrame = [{
    label: "All",
    value: ""
  }];
  const service = new Services(datasource.meta.jsonData.cmdbEndpoint || "", datasource.meta.jsonData.grafanaEndpoint || "");
  const [elementId, setElementId] = useState("");
  const [supportedPanels, setSupportedPanels] = useState([]);
  const [allFrames, setAllFrames] = useState<Record<string, string>>({});
  const [frames, setFrames] = useState(defaultFrame);
  const onChanged = useRef(false);

  const getCloudElements = useCallback((id: string, query: any) => {
    service.getCloudElements(id).then((res: any) => {
      if (res && res[0]) {
        const cloudElement = res[0];
        if (cloudElement) {
          const cloudElementQuery = getCloudElementsQuery(id, cloudElement, datasource.meta.jsonData.awsxEndPoint || "");
          query = {
            ...query,
            ...cloudElementQuery
          };
          onChange({ ...query });
          service.getSupportedPanels(cloudElement.elementType.toUpperCase(), "AWS").then((res) => {
            if (res && res.length > 0) {
              const panels: any = [];
              const frames: any = {};
              res.map((panel: any) => {
                panels.push({
                  label: panel.name,
                  value: panel.name
                });
                frames[panel.name] = panel.frames;
              });
              setSupportedPanels(panels);
              setAllFrames(frames);
              const { queryString } = query;
              const frame = frames[queryString];
              if (frame) {
                const arrFrames: any = frame.split(",").map((f: any) => {
                  return {
                    label: f,
                    value: f
                  };
                });
                setFrames([...defaultFrame, ...arrFrames]);
              } else {
                setFrames([...defaultFrame]);
              }
            }
          });
        }
      }
    });
  }, [onChange]); // eslint-disable-line react-hooks/exhaustive-deps

  useEffect(() => {
    if (onChanged.current === false) {
      let id = "";
      if (document.getElementById("elementId")) {
        id = (document.getElementById("elementId") as HTMLInputElement)?.value;
      } else if (document.getElementById("var-elementId")) {
        id = (document.getElementById("var-elementId") as HTMLInputElement)?.value;
      } else {
        id = query.elementId;
      }
      if (id) {
        setElementId(id);
        getCloudElements(id, query);
      } else {
        alert("Please set 'elementId' variable");
      }
      onChanged.current = true;
    }
  }, [query, onChange, getCloudElements]);

  const onChangeElementType = (e: any) => {
    onChange({ ...query, elementType: e.target.value });
  };

  const onChangeInstanceID = (e: any) => {
    onChange({ ...query, cloudIdentifierId: e.target.value });
  };

  const onChangeSupportedPanel = (value: any) => {
    onChange({ ...query, queryString: value });
    const frame = allFrames[value];
    if (frame) {
      const arrFrames: any = frame.split(",").map(f => {
        return {
          label: f,
          value: f
        };
      });
      setFrames([...defaultFrame, ...arrFrames]);
    } else {
      setFrames([...defaultFrame]);
    }
  };

  const onChangeResponseType = (value: any) => {
    onChange({ ...query, responseType: value });
  };

  const onChangeFrame = (value: any) => {
    localStorage.setItem("datasource-selected-frame", value);
    onChange({ ...query, selectedFrame: value });
  };

  const {
    elementType,
    cloudIdentifierId,
    queryString,
    responseType,
    selectedFrame
  } = query;

  return (
    <div>
      <EditorRows>
        <EditorRow label="">
          <InlineField label="Element Type">
            <Input value={elementType} onChange={(e: any) => onChangeElementType(e)} />
          </InlineField>
          <InlineField label="Instance ID">
            <Input value={cloudIdentifierId} onChange={(e: any) => onChangeInstanceID(e)} />
          </InlineField>
          <InlineField label="Element ID">
            <Input disabled={true} value={elementId} />
          </InlineField>
        </EditorRow>
      </EditorRows>
      <EditorRows>
        <EditorRow label="">
          <InlineField label="Supported Panels">
            <Select
              className="min-width-12 width-12"
              value={queryString}
              options={supportedPanels}
              onChange={(e) => onChangeSupportedPanel(e.value)}
              menuShouldPortal={true}
            />
          </InlineField>
          <InlineField label="Response Type">
            <Select
              className="min-width-12 width-12"
              value={responseType}
              options={RESPONSE_TYPE}
              onChange={(e) => onChangeResponseType(e.value)}
              menuShouldPortal={true}
            />
          </InlineField>
          <InlineField label="Frames">
            <Select
              className="min-width-12 width-12"
              value={selectedFrame}
              options={frames}
              onChange={(e) => onChangeFrame(e.value)}
              menuShouldPortal={true}
            />
          </InlineField>
        </EditorRow>
      </EditorRows>
    </div>
  );
}
