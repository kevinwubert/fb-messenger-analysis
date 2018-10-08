import React from "react";

import './Graph.css';

class Graph extends React.Component {
  render() {
    return (
      <div className="Graph">
        <h1> result </h1>
        <img src={this.props.url}></img>
      </div>
    );
  }
}
export default Graph;