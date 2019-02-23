import React, { Component } from 'react';
import styled, { css } from 'styled-components'

const AppContainer = styled.div`
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100vh;
`;

const Title = styled.h1`
  text-align: center;
`;

const Input = styled.input`
  border: 1px solid lightgrey;
  border-radius: 8px;
  padding: 8px;
  font-size: 12pt;
  width: 250px;
`;

const Form = styled.form``;

const FormElement = styled.div`
  margin: 16px;
  flex: 1;
`;

const Button = styled.button`
  width: 100%;
  padding: 8px;
  border-radius: 8px;
  background: lightblue;
  color: white;
  font-size: 12pt;
  border: 2px solid lightblue;

  ${props => props.secondary &&
    css`
      background: white;
      color: lightblue;
    `};
`;

class App extends Component {
  render() {
    return (
      <AppContainer>
        <Form>
          <Title>HalsPals</Title>

          <FormElement>
            <Input id="email" type="email" placeholder="john@email.com"/>
          </FormElement>

          <FormElement>
            <Input id="password" type="password" placeholder="password"/>
          </FormElement>
          
          <FormElement>
            <Button>Login</Button>
          </FormElement>
          
          <FormElement>
            <Button secondary>Register</Button>
          </FormElement>
        </Form>
      </AppContainer>
    );
  }
}

export default App;
