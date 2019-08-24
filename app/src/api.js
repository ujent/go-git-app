import queryString from 'query-string';

export const apiUri = (function () {
    const buildMode = process.env.REACT_APP_BUILD_MODE;
    if (!buildMode) return 'http://localhost:3000';
    switch (buildMode) {
        case 'dev':
            return 'http://localhost:3000';
        case 'staging':
            return 'http://localhost:3000';
        case 'release':
            return 'http://localhost:3000';
        default:
            return '';
    }
})();

function fetchApi(url, options = {}) {
    const defaulOptions = {
        /*credentials: 'include'*/
    };

    return fetch(apiUri + url, Object.assign(defaulOptions, options)).then(
        response => {
            if (response.ok) {
                const res = response.json();
                return res;
            } else {
                console.log(response);
                return response.text().then(err => {
                    return Promise.reject({
                        message: err,
                        status: response.status
                    });
                });
            }
        }
    );
}

export function switchUser(name) {
    const query = JSON.stringify({ name: name });

    return fetchApi('/users/switch', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: query
    });
}

export function getRepositories() {
    return fetchApi('/repositories', {});
}

export function switchRepo(name) {

    return fetchApi('/repositories/open/' + name, {
    });
}

export function getBranches() {
    return fetchApi('/branches', {});
}

export function switchBranch(name) {
    return fetchApi('/branches/checkout/' + name, {});
}

export function getRepoFiles() {
    return fetchApi('files', {});
}

export function commit() {
}
export function checkoutBranch() {
}
export function clone() {
}
export function log() {
}
export function merge() {
}
export function pull() {
}
export function push() {
}
export function removeBranch() {
}
export function removeRepo() {
}
