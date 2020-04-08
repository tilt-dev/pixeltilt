import { useState } from "react";
import Header from "../components/Header";
import styled from "styled-components";
import fetch from "isomorphic-unfetch";
import axios from "axios";

let MainPane = styled.main`
  display: flex;
  justify-content: center;
  align-items: center;
  width: 100vw;
  height: 100vh;
`;

let MaxImg = styled.img`
  max-width: 60vw;
  max-height: 60vh;
`;

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

  const handleUploadFile = e => {
    const file = e.currentTarget.files[0];
    setFileSelection(file);
  };

  const apply = async () => {
    const data = new FormData();
    const file = fileSelection;
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
      return <MaxImg src={`/api/image/${resultingImage}`}></MaxImg>;
    }

    if (fileSelection) {
      let url = URL.createObjectURL(fileSelection);
      return <MaxImg src={url}></MaxImg>;
    }

    return (
      <div>
        <div>
          <label htmlFor="upload">File to upload</label>
        </div>
        <div>
          <input type="file" onChange={handleUploadFile} />
        </div>
      </div>
    );
  };

  return (
    <div>
      <Header
        filters={filters}
        clearAndSetCheckedItems={clearAndSetCheckedItems}
        checkedItems={checkedItems}
        reset={reset}
        apply={apply}
        hasFileSelection={!!fileSelection}
        statusMessage={statusMessage()}
      />
      <MainPane>{renderContent()}</MainPane>
    </div>
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
