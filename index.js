'use strict';
const e = React.createElement;
class Button extends React.Component {
  constructor(props) {
    super(props);
    this.state = { clicked: false };
  }
  render() {
    if (this.state.clicked) {
      return 'Thanks for clicking!';
    }
  return e(
      'button',
      { onClick: () => this.setState({ clicked: true }) },
      'Submit'
    );
  }
}
const domContainer = document.querySelector('#button_container');
ReactDOM.render(e(Button), domContainer);