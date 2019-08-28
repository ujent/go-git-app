import React, { Component } from 'react';
import '../../App.css';

export default class Clone extends Component {
    constructor(props) {
        super(props);

        this.state = {
            url: '',
            name: '',
            psw: ''
        }
    }

    render() {
        const isAvailable = this.props.isAvailable && this.state.url;

        const onURLChange = (e) => {
            this.setState({
                url: e.target.value
            });
        }

        const onAuthNameChange = (e) => {
            this.setState({
                name: e.target.value
            });
        }

        const onAuthPswChange = (e) => {
            this.setState({
                psw: e.target.value
            });
        }

        const onCloneClick = () => {
            this.props.action(this.state.url, this.state.name, this.state.psw);
            this.setState({
                url: '',
                name: '',
                psw: ''
            });
        }

        return (
            < li >
                <div className="command-block command-block-content">
                    <button type="button" className="button" disabled={!isAvailable} onClick={() => onCloneClick()}>Clone</button>
                    <input type="text" placeholder="URL" className="command-block-input" value={this.state.url} onChange={(e) => onURLChange(e)} />
                </div>
                <div className="credentials">
                    <input type="text" placeholder="user" value={this.state.name} onChange={(e) => onAuthNameChange(e)} />
                    <input type="password" placeholder="password" value={this.state.psw} onChange={(e) => onAuthPswChange(e)} />
                </div>
                <hr></hr>
            </li >
        );

    };
};
