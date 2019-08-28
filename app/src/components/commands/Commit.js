import React, { Component } from 'react';
import '../../App.css';

export default class Commit extends Component {
    constructor(props) {
        super(props);

        this.state = {
            comment: ''
        }
    }

    render() {
        const onCommentChange = (e) => {
            this.setState({
                comment: e.target.value
            });
        }

        const onCommitClick = () => {
            this.props.action(this.state.comment);

            this.setState({
                comment: ''
            })
        }

        return (
            <li className="commit-command">
                <button type="button" className="button" disabled={!this.props.isAvailable} onClick={() => onCommitClick()}>Commit</button>
                <textarea placeholder="commit message" rows="3" value={this.state.comment} onChange={(e) => onCommentChange(e)}></textarea>
                <hr></hr>
            </li>
        );
    }

}
