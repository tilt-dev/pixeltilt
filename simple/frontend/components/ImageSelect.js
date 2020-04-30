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
  let img = props.img;
  let selectImage = props.selectImage;

  let onClick = e => {
    selectImage(img);
  };
  return <Img src={img.url} onClick={onClick} />;
}

export default ImageSelect;
