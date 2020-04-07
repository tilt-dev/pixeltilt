import styled from "styled-components";

let HeaderRoot = styled.header`
  position: absolute;
  width: 100%;
  color: #586e75;
  box-sizing: border-box;
`;

let HeaderControls = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
  padding: 0 72px;
  min-height: 80px;
  background-color: #002b36;
  box-sizing: border-box;
`;

let HeaderStatus = styled.div`
  display: flex;
  justify-content: center;
  align-items: center;
  width: 100%;
  min-height: 40px;
  background-color: #c4c4c4;
  box-sizing: border-box;
`;

let HeaderMiddleCell = styled.div``;

const Header = props => {
  const filters = props.filters;
  const reset = props.reset;
  const apply = props.apply;
  const setCheckedItems = props.setCheckedItems;
  const checkedItems = props.checkedItems;
  const statusMessage = props.statusMessage;

  const handleChangeFilter = event => {
    setCheckedItems({
      ...checkedItems,
      [event.target.name]: event.target.checked
    });
  };

  let filterEls = filters.map(item => {
    return (
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
    );
  });

  return (
    <HeaderRoot>
      <HeaderControls>
        <div>Pixel Tilt!</div>
        <HeaderMiddleCell>
          {filterEls}
          <button onClick={apply}>Apply Tilt</button>
        </HeaderMiddleCell>
        <button onClick={reset}>Reset button</button>
      </HeaderControls>
      <HeaderStatus>{statusMessage}</HeaderStatus>
    </HeaderRoot>
  );
};

export default Header;
