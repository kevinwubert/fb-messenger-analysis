import React, { Component } from 'react';
import Graph from "./Graph";
import './App.css';

import { Jumbotron, Col, Button, Form, FormGroup, FormControl, ControlLabel } from 'react-bootstrap';

const url = 'http://35.236.58.233/';

class App extends Component {
  constructor(props) {
    super(props);

    this.state = {
      names: [],
      url: ''
    }

    this.handleChange = this.handleChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
  }
  
  componentDidMount() {
    fetch(url + 'getNames')
      .then(response => response.json())
      .then(data => this.setState({ names: data }));
  }

  handleChange(e) {
    var newState = Object.assign({}, this.state);
    newState.
    this.setState({ value: e.target.value });
  }


  handleSubmit(event) {
    event.preventDefault();
    this.setState({
      names: this.state.names,
      url: `http://35.236.58.233/graph?name=${this.name.value}&type=${this.type.value}&count=${this.count.value}`
    })
  }
  
  render() {
    var types = ['words', 'stickers', 'mentions', 'reactions'];
    return (
      <div className="App">
        <Jumbotron>
          <h1>Facebook Messenger Analysis</h1>
          <p>
            Simple app to view frequency/trends of chat trends on Facebook Messenger
          </p>
        </Jumbotron>

        <Form horizontal onSubmit={this.handleSubmit}>
          <FormGroup controlId="name">
            <Col componentClass={ControlLabel} sm={2}>
              Name
            </Col>
            <Col sm={8}>
              <FormControl
                componentClass="select" 
                inputRef={ref => {this.name = ref; }}
                placeholder="Everyone">
                {
                  this.state.names.map((name, index) => {
                    return <option value={name}>{name}</option>
                  })
                }
              </FormControl>
            </Col>
          </FormGroup>

          <FormGroup controlId="type">
            <Col componentClass={ControlLabel} sm={2}>
              Types
            </Col>
            <Col sm={8}>
              <FormControl 
                componentClass="select" 
                inputRef={ref => { this.type = ref; }}
                placeholder="Words">
                {
                  types.map((type, index) => {
                    return <option value={type}>{type}</option>
                  })
                }
              </FormControl>
            </Col>
          </FormGroup>

          <FormGroup controlId="count">
            <Col componentClass={ControlLabel} sm={2}>
              Count
            </Col>
            <Col sm={8}>
                <FormControl 
                  type="number" 
                  inputRef={ref => { this.count = ref; }}
                  defaultValue = "20"
                  placeholder="Enter number of words to return" />
            </Col>
          </FormGroup>

          <Button type="submit">Submit</Button>

        </Form>

        <Graph url={this.state.url}></Graph>
      </div>
    );
  }
}

export default App;
