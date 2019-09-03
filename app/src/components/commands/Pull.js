import React, { Component } from 'react';
import '../../App.css';

export default class Pull extends Component {
    constructor(props) {
        super(props);

        this.state = {
            url: '',
            name: '',
            psw: ''
        }
    }

    render() {

        return (
            <li>
                <div className="command-block command-block-content">
                    <button type="button" className="button" disabled={!this.props.isAvailable} onClick={() => this.onPullClick()}>Pull</button>
                    <input type="text" placeholder="remote" className="command-block-input" value={this.state.url} onChange={(e) => this.onURLChange(e)} />
                </div>
                <div className="credentials">
                    <input type="text" placeholder="user" value={this.state.name} onChange={(e) => this.onAuthNameChange(e)} />
                    <input type="password" placeholder="password" value={this.state.psw} onChange={this.onAuthPswChange} />
                </div>
                <hr></hr>
            </li>
        );
    };

    onURLChange = (e) => {
        this.setState({
            url: e.target.value
        });
    }

    onAuthNameChange = (e) => {
        this.setState({
            name: e.target.value
        });
    }

    onAuthPswChange = (e) => {
        this.setState({
            psw: e.target.value
        });
    }

    onPullClick = () => {
        this.props.action(this.state.url, this.state.name, this.state.psw);
        this.setState({
            url: '',
            name: '',
            psw: ''
        });
    }

}
