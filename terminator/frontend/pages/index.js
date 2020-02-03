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
      ...defaultCheckedItems,
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
    data.append("filters", JSON.stringify(checkedItems));
    let response = await axios.post("/api/upload", data, {});
    console.log(response.data);
    setResultingImage(response.data);
  };

  // TODO(dmiller): this should display the resulting image from the returned image URL
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
              name={item.label}
              onChange={handleChangeFilter}
              checked={checkedItems[item.label]}
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
    </div>
  );
};

Index.getInitialProps = async function() {
  // TODO(dmiller): we could await on these things more efficiently
  const filtersRes = await fetch("http://muxer/filters");
  const filtersData = await filtersRes.json();

  const imagesRes = await fetch("http://muxer/images");
  const imagesData = await imagesRes.json();

  const defaultCheckedItems = {};
  filtersData.forEach(c => (defaultCheckedItems[c.label] = true));

  return {
    filters: filtersData,
    images: imagesData,
    defaultCheckedItems: defaultCheckedItems
  };
};

export default Index;
