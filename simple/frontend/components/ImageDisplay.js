import styled, { keyframes } from "styled-components";
import { ReactComponent as ImageMarkSvg } from "./tilt-imagemark.svg";

const rotate = keyframes`
  from {
    transform: rotate(0deg);
  }

  to {
    transform: rotate(360deg);
  }
`;

let Container = styled.div`
  max-width: 60vw;
  max-height: 60vh;
  position: relative;
`;

let Overlay = styled.div`
  position: absolute;
  align-items: center;
  justify-content: center;
  top: 0;
  bottom: 0;
  right: 0;
  left: 0;
  display: none;

  ${Container}.is-pending & {
    display: flex;
  }
`;

let SpinningImageMark = styled(ImageMarkSvg)`
  animation: ${rotate} 2s linear infinite;
  width: 100px;
  height: 100px;
  opacity: 0;
  transition: opacity 300ms ease;

  ${Container}.is-pending & {
    opacity: 1;
  }
`;

let Img = styled.img`
  max-width: 60vw;
  max-height: 60vh;
  transition: opacity 300ms ease;

  ${Container}.is-pending & {
    opacity: 0.64;
  }
`;

function ImageDisplay(props) {
  let extraClasses = [];
  if (props.isPending) {
    extraClasses.push("is-pending");
  }

  return (
    <Container className={extraClasses.join(" ")}>
      <Img src={props.src} />
      <Overlay>
        <SpinningImageMark />
      </Overlay>
    </Container>
  );
}

export default ImageDisplay;
