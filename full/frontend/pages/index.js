import { useState } from "react";
import Header from "../components/Header";
import styled from "styled-components";
import fetch from "isomorphic-unfetch";
import axios from "axios";
import UploadControl from "../components/UploadControl"
import ImageSelect from "../components/ImageSelect"
import ImageDisplay from "../components/ImageDisplay"
const babyBear = "/baby-bear.png"
const plane = "/plane.png"
const duck = "/duck.png"

let Root = styled.div`
  width: 100%;
  min-height: 100vh;
  display: flex;
  flex-direction: column;
`

let MainPane = styled.main`
  display: flex;
  justify-content: center;
  align-items: center;
  flex-grow: 1;
`;

let ImageGrid = styled.div`
  display: grid;
  grid-template-columns: 50% 50%;
  grid-template-rows: auto auto;
`

const Index = props => {
  const filters = props.filters;
  const defaultCheckedItems = props.defaultCheckedItems;

  const [checkedItems, setCheckedItems] = useState(defaultCheckedItems);
  const [fileSelection, setFileSelection] = useState();
  const [resultingImage, setResultingImage] = useState("");
  const [applyState, setApplyState] = useState({
    inProgress: false,
    error: null
  });

  const setFileBlob = (blob) => {
    let url = URL.createObjectURL(blob);
    setFileSelection({blob, url});
  }

  const statusMessage = () => {
    if (applyState.inProgress) {
      return "Applying filter…";
    }

    if (applyState.error) {
      return "Error applying filter: " + applyState.error;
    }

    if (resultingImage || fileSelection) {
      return "Select one or more filters and apply";
    }

    return "Upload an image. Apply filters once selected";
  };

  const clearAndSetCheckedItems = checkedItems => {
    setResultingImage("");
    setCheckedItems(checkedItems);
    setApplyState({ inProgress: false, error: null });
  };

  const reset = () => {
    setResultingImage("");
    setFileSelection(null);
    setCheckedItems(defaultCheckedItems);
    setApplyState({ inProgress: false, error: null });
  };

  const apply = async () => {
    if (!fileSelection) {
      throw new Error('internal error: no file to apply filters on')
    }

    const data = new FormData();
    const file = fileSelection.blob;
    data.append("file", file);
    const filters = Object.keys(checkedItems)
      .filter(key => checkedItems[key])
      .reduce((res, key) => ((res[key] = checkedItems[key]), res), {});
    data.append("filters", JSON.stringify(filters));
    setApplyState({ inProgress: true });

    axios
      .post("/api/upload", data, {})
      .then(response => {
        setResultingImage(response.data.name);
        setApplyState({ inProgress: false });
      })
      .catch(err => {
        setApplyState({ error: err });
      });
  };

  const renderContent = () => {
    if (resultingImage) {
      return <ImageDisplay src={`/api/image/${resultingImage}`} />;
    }

    if (fileSelection) {
      return <ImageDisplay src={fileSelection.url} isPending={applyState.inProgress} />;
    }

    return (
      <ImageGrid>
        <UploadControl setFileBlob={setFileBlob} />
        <ImageSelect url={babyBear} setFileBlob={setFileBlob} />
        <ImageSelect url={plane} setFileBlob={setFileBlob} />
        <ImageSelect url={duck} setFileBlob={setFileBlob} />
      </ImageGrid>
    );
  };

  return (
    <Root>
      <Header
        filters={filters}
        clearAndSetCheckedItems={clearAndSetCheckedItems}
        checkedItems={checkedItems}
        reset={reset}
        apply={apply}
        hasFileSelection={!!fileSelection}
        statusMessage={statusMessage()}
      />
      <MainPane>
        {renderContent()}
      </MainPane>
    </Root>
  );
};

Index.getInitialProps = async function() {
  const filtersRes = await fetch("http://muxer:8080/filters");
  const filtersData = await filtersRes.json();

  const imagesRes = await fetch("http://muxer:8080/images");
  const imagesData = await imagesRes.json();

  const defaultCheckedItems = {};
  filtersData.forEach(
    c => (defaultCheckedItems["filter_" + c.label.toLowerCase()] = false)
  );

  return {
    filters: filtersData,
    images: imagesData,
    defaultCheckedItems: defaultCheckedItems
  };
};

export default Index;
