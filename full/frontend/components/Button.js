import styled from "styled-components";
import color from "./color";

let ButtonRoot = styled.button`
  display: flex;
  border: 1px solid ${color.red};
  color: ${color.red};
  box-sizing: border-box;
  border-radius: 4px;
  padding: 0 8px;
  margin: 0 8px 0 0;
  align-items: center;
  justify-content: space-between;
  height: 30px;
  background-color: inherit;
  font-size: 10px;
  line-height: 14px;
  font-family: inherit;
  cursor: pointer;
  transition: background-color 300ms ease, color 300ms ease;
  box-sizing: content-box;

  &.is-toggle:hover {
    padding: 1px 9px;
    margin-left: -1px;
    margin-right: 7px !important;
  }

  &[disabled],
  &.is-green[disabled],
  &.is-yellow[disabled] {
    pointer-events: none;
    border-color: ${color.gray};
    color: ${color.gray};
  }
  &[disabled] path,
  &.is-green[disabled] path,
  &.is-yellow[disabled] path {
    fill: ${color.gray};
  }

  & > svg {
    width: 14px;
    height: 14px;
    margin-right: 7px;
  }

  & path {
    transition: fill 300ms ease;
    fill: ${color.red};
  }

  &:not(.is-toggle):hover {
    background-color: ${color.red};
    color: ${color.grayDark};
  }
  &:not(.is-toggle):hover path {
    fill: ${color.grayDark};
  }

  &.is-yellow {
    border-color: ${color.yellow};
    color: ${color.yellow};
  }
  &.is-yellow path {
    fill: ${color.yellow};
  }
  &.is-yellow:hover {
    color: ${color.grayDark};
    background-color: ${color.yellow};
  }
  &.is-yellow:hover path {
    fill: ${color.grayDark};
  }

  &.is-green {
    border-color: ${color.green};
    color: ${color.green};
  }
  &.is-green path {
    fill: ${color.green};
  }
  &.is-green:hover {
    color: ${color.grayDark};
    background-color: ${color.green};
  }
  &.is-green:hover path {
    fill: ${color.grayDark};
  }

  &.is-selected {
    background-color: ${color.red};
    color: ${color.grayDark};
  }
  &.is-selected.is-yellow {
    background-color: ${color.yellow};
  }
  &.is-selected.is-green {
    background-color: ${color.green};
  }
  &.is-selected path {
    fill: ${color.grayDark};
  }
`;

function Button(props) {
  let onClick = props.onClick;
  let name = props.name;
  let extraClasses = [];
  if (props.selected) {
    extraClasses.push("is-selected");
  }
  if (props.isToggle) {
    extraClasses.push("is-toggle");
  }
  if (props.isGreen) {
    extraClasses.push("is-green");
  }
  if (props.isYellow) {
    extraClasses.push("is-yellow");
  }
  let disabled = props.disabled;
  return (
    <ButtonRoot
      onClick={onClick}
      name={name}
      style={props.style}
      disabled={disabled}
      className={extraClasses.join(" ")}
    >
      {props.children}
    </ButtonRoot>
  );
}

export default Button;
