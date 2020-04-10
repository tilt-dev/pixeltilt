import styled from "styled-components";
import color from "./color";

let UploadControlRoot = styled.div`
  width: 220px;
  height: 220px;
  margin: 24px;
  border: 2px dashed ${color.grayLighter};
  border-radius: 6px;
  box-sizing: border-box;
  transition: background-color 300ms ease;

  &:hover {
    border: 2px solid ${color.grayLighter};
    background-color: ${color.grayLighter};
  }
`;

let UploadControlLabel = styled.label`
  display: flex;
  width: 100%;
  height: 100%;
  justify-content: center;
  align-items: center;
  font-size: 20px;
  font-weight: 700;
  line-height: 27px;
  color: ${color.gray};
  cursor: pointer;
  padding: 24px;
  box-sizing: border-box;
  text-align: center;
`;

function UploadControl(props) {
  let setFileBlob = props.setFileBlob;

  const handleUploadFile = e => {
    const file = e.currentTarget.files[0];
    setFileBlob(file);
  };

  return (
    <UploadControlRoot>
      <input
        type="file"
        id="upload-file"
        onChange={handleUploadFile}
        style={{ display: "none" }}
      />
      <UploadControlLabel htmlFor="upload-file">
        Click to upload a PNG
      </UploadControlLabel>
    </UploadControlRoot>
  );
}

export default UploadControl;
