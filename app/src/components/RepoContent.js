import React, { Component } from 'react';
import '../App.css';

export default class RepoContent extends Component {
    constructor(props) {
        super(props);

        this.state = {
            selectedFile: '',
        }
    }

    render() {
        const onFileClick = (name) => {
            this.setState({
                selectedFile: name
            })
        }

        const files = this.props.files.map(el => {
            if (el.path === this.state.selectedFile) {
                if (el.isConflict) {
                    return <li className="repo-file repo-file-selected conflict-file" onClick={() => onFileClick(el.path)}>{el.path}</li>
                }

                return <li className="repo-file repo-file-selected" onClick={() => onFileClick(el.path)}>{el.path}</li>
            }

            if (el.isConflict) {
                return <li className="repo-file conflict-file" onClick={() => onFileClick(el.path)}>{el.path}</li>
            }

            return <li className="repo-file" onClick={() => onFileClick(el.path)}>{el.path}</li>
        });

        return (
            <section className="repo-content">
                <h2 className="visually-hidden">Repository content</h2>
                <div className="repo-files">
                    <h3>Files</h3>
                    <ul className="files-list">
                        {files}
                    </ul>
                </div>
                {
                    this.state.selectedFile ?
                        <div className="file-content ">
                            <h3>{this.state.selectedFile}</h3>
                            <div className="file-content-buttons">
                                <button type="button" className="button">Save</button>
                            </div>
                            <textarea rows="32"></textarea>
                        </div> : null
                }

            </section>
        );
    };

}