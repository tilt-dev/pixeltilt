import styled from "styled-components";
import { ReactComponent as PixelTiltLogoSvg } from "./pixel-tilt-logo.svg";
import { ReactComponent as TiltLogoSvg } from "./tilt-logo.svg";
import { ReactComponent as DropSvg } from "./drop.svg";
import { ReactComponent as MixerSvg } from "./mixer.svg";
import { ReactComponent as ScanSvg } from "./scan.svg";
import { ReactComponent as ResetSvg } from "./reset.svg";
import { ReactComponent as ApplySvg } from "./apply.svg";
import Button from "./Button";
import color from "./color";

let HeaderRoot = styled.header`
  color: ${color.gray};
  box-sizing: border-box;
`;

let HeaderControls = styled.div`
  display: flex;
  align-items: center;
  min-height: 80px;
  background-color: ${color.grayDark};
  box-sizing: border-box;
  justify-content: space-between;
`;

let HeaderPixelLogoCell = styled.a`
  display: flex;
  padding-left: 72px;
  width: 240px;
  box-sizing: border-box;
`;

let HeaderTiltLogoCell = styled.a`
  padding-right: 72px;
  display: flex;
  flex-direction: row-reverse;
  width: 240px;
  box-sizing: border-box;
  color: ${color.gray} !important;
  text-decoration: none !important;
`;

let TiltBrandText = styled.div`
  font-family: "Fira Code", sans-serif;
  font-size: 9px;
  line-height: 13px;
  text-align: right;
  width: 80px;
`;

let HeaderStatus = styled.div`
  display: flex;
  justify-content: center;
  align-items: center;
  width: 100%;
  min-height: 40px;
  background-color: ${color.grayLight};
  box-sizing: border-box;
`;

let HeaderStatusMessage = styled.div`
  max-width: 650px;
  text-align: left;
  flex-grow: 1;
`;

let HeaderButtonRow = styled.div`
  display: flex;
  max-width: 650px;
  justify-content: space-between;
  align-items: center;
  height: 40px;
  flex-grow: 3;
`;

let HeaderButtonInnerRow = styled.div`
  display: flex;
  align-items: center;
  height: 40px;
`;

let LogoLink = styled.a`
  display: flex;
`;

const Header = props => {
  const filters = props.filters;
  const reset = props.reset;
  const apply = props.apply;
  const clearAndSetCheckedItems = props.clearAndSetCheckedItems;
  const checkedItems = props.checkedItems;
  const statusMessage = props.statusMessage;
  const hasFileSelection = props.hasFileSelection;
  const hasFilterSelection = Object.values(checkedItems).some(v => v);

  const handleChangeFilter = event => {
    let target = event.currentTarget;
    let newCheckedItems = { ...checkedItems };
    newCheckedItems[target.name] = !checkedItems[target.name];
    clearAndSetCheckedItems(newCheckedItems);
  };

  let filterEls = filters.map(item => {
    let name = "filter_" + item.label.toLowerCase();
    let selected = checkedItems[name];
    let label = item.label;
    let svg = null;
    if (item.label == "Red") {
      label = "Color";
      svg = <DropSvg />;
    } else if (item.label == "Glitch") {
      svg = <MixerSvg />;
    } else if (item.label == "Rectangler") {
      label = "Object";
      svg = <ScanSvg />;
    }
    return (
      <Button
        onClick={handleChangeFilter}
        name={name}
        key={name}
        selected={selected}
        isToggle={true}
        disabled={!hasFileSelection}
      >
        {svg}
        <div>{label}</div>
      </Button>
    );
  });

  return (
    <HeaderRoot>
      <HeaderControls>
        <HeaderPixelLogoCell href="/">
          <PixelTiltLogoSvg />
        </HeaderPixelLogoCell>
        <HeaderButtonRow>
          <HeaderButtonInnerRow>
            {filterEls}

            <Button
              onClick={apply}
              style={{ marginLeft: "16px" }}
              isGreen={true}
              disabled={!hasFileSelection || !hasFilterSelection}
            >
              <ApplySvg />
              <div>Apply Tilt</div>
            </Button>
          </HeaderButtonInnerRow>

          <Button onClick={reset} isYellow={true}>
            <ResetSvg />
            <div>Reset</div>
          </Button>
        </HeaderButtonRow>

        <HeaderTiltLogoCell href="https://tilt.dev/">
          <TiltLogoSvg height="26px" />
          <TiltBrandText>A Sample App showcasing!</TiltBrandText>
        </HeaderTiltLogoCell>
      </HeaderControls>
      <HeaderStatus>
        <HeaderStatusMessage>{statusMessage}</HeaderStatusMessage>
      </HeaderStatus>
    </HeaderRoot>
  );
};

export default Header;
