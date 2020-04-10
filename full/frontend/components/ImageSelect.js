import styled from "styled-components";
import fetch from "isomorphic-unfetch";

let Img = styled.img`
  width: 220px;
  height: 220px;
  margin: 24px;
  object-fit: cover;
  cursor: pointer;

  &:hover {
    width: 236px;
    height: 236px;
    margin: 16px;
  }
`;

function ImageSelect(props) {
  let url = props.url;
  let setFileBlob = props.setFileBlob;

  let onClick = e => {
    fetch(url)
      .then(response => response.blob())
      .then(blob => setFileBlob(blob))
      .catch(err => {
        console.error("Error converting URL into blob", err);
      });
  };
  return <Img src={url} onClick={onClick} />;
}

export default ImageSelect;
