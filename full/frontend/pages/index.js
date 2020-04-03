import { useState } from "react";
import fetch from "isomorphic-unfetch";
import axios from "axios";
import btoa from "btoa";

const Index = props => {
  const filters = props.filters;
  const defaultCheckedItems = props.defaultCheckedItems;

  const [checkedItems, setCheckedItems] = useState(defaultCheckedItems);
  const [fileToUpload, setFileToUpload] = useState();
  const [resultingImage, setResultingImage] = useState("");

  const handleChangeFilter = event => {
    setCheckedItems({
      ...checkedItems,
      [event.target.name]: event.target.checked
    });
  };

  const handleUploadFile = e => {
    const file = e.currentTarget.files[0];
    setFileToUpload(file);
  };

  const uploadFile = async () => {
    const data = new FormData();
    const file = fileToUpload;
    data.append("file", file);
    const filters = Object.keys(checkedItems)
      .filter(key => checkedItems[key])
      .reduce((res, key) => ((res[key] = checkedItems[key]), res), {});
    data.append("filters", JSON.stringify(filters));
    let response = await axios.post("/api/upload", data, {});
    setResultingImage(response.data.name);
  };

  return (
    <div>
      <label>Checked item name : {JSON.stringify(props.filtersData)}</label>{" "}
      <br />
      <form>
        {filters.map(item => (
          <label key={item.label}>
            {item.label}
            <input
              type="checkbox"
              id={item.label}
              name={"filter_" + item.label.toLowerCase()}
              onChange={handleChangeFilter}
              checked={checkedItems["filter_" + item.label.toLowerCase()]}
            />
          </label>
        ))}
        <br />
        <label>
          <br />
          File to upload : {JSON.stringify(fileToUpload)}
          <br />
          <input type="file" onChange={handleUploadFile} />
        </label>
        <button type="button" onClick={uploadFile}>
          Render
        </button>
      </form>
      {resultingImage !== "" && (
        <img src={`/api/image/${resultingImage}`}></img>
      )}
    </div>
  );
};

Index.getInitialProps = async function() {
  // TODO(dmiller): we could await on these things more efficiently
  const filtersRes = await fetch("http://muxer:8080/filters");
  const filtersData = await filtersRes.json();

  const imagesRes = await fetch("http://muxer:8080/images");
  const imagesData = await imagesRes.json();

  const defaultCheckedItems = {};
  filtersData.forEach(
    c => (defaultCheckedItems["filter_" + c.label.toLowerCase()] = true)
  );

  return {
    filters: filtersData,
    images: imagesData,
    defaultCheckedItems: defaultCheckedItems
  };
};

export default Index;
