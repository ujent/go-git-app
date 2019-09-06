import queryString from 'query-string';

export const apiUri = (function () {
    const buildMode = process.env.REACT_APP_BUILD_MODE;
    if (!buildMode) return 'http://localhost:4000';
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

export function getRepositories(user) {
    return fetchApi('/repositories/' + user, {});
}

export function getCurrentRepository(user) {
    return fetchApi('/repositories/current/' + user, {});
}

export function switchRepo(user, repo) {
    const query = queryString.stringify({ repo: repo, user: user });

    return fetchApi('/repositories/open?' + query, {});
}


export function getBranches(user, repo) {
    const query = queryString.stringify({ repo: repo, user: user });

    return fetchApi('/branches?' + query, {});
}

export function checkoutBranch(user, repo, branch) {
    const query = JSON.stringify({ repo: repo, user: user, branch: branch });

    return fetchApi('/branches/checkout', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: query
    });
}

export function getFiles(user, repo, branch) {
    const query = queryString.stringify({ repo: repo, user: user, branch: branch });

    return fetchApi('/files/all?' + query, {});
}

export function getFile(settings, path, isConflict) {
    const query = queryString.stringify({ branch: settings.branch, repo: settings.repo, user: settings.user, path: path, isConflict: isConflict });

    return fetchApi('/files?' + query, {});
}

export function addFile(settings, path, content) {
    const query = JSON.stringify({ base: settings, path: path, content: content });

    return fetchApi('/files', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: query
    });
}

export function editFile(settings, path, content, isConflict) {
    const query = JSON.stringify({ base: settings, path: path, content: content, isConflict: isConflict });

    return fetchApi('/files', {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json'
        },
        body: query
    });
}

export function removeFile(settings, path, isConflict) {
    const query = JSON.stringify({ base: settings, path: path, isConflict: isConflict });

    return fetchApi('/files', {
        method: 'DELETE',
        headers: {
            'Content-Type': 'application/json'
        },
        body: query
    });
}

export function commit(settings, msg) {
    const query = JSON.stringify({ base: settings, message: msg });

    return fetchApi('/commit', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: query
    });
}

export function clone(user, url, authName, authPsw) {
    const query = authName ?
        JSON.stringify({ user: user, URL: url, auth: { name: authName, psw: authPsw } })
        : JSON.stringify({ user: user, URL: url });

    return fetchApi('/repositories/clone', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: query
    });
}


export function log(user, repo, branch) {
    const query = queryString.stringify({ repo: repo, user: user, branch: branch });

    return fetchApi('/log?' + query, {});
}

export function pull(settings, remote, authName, authPsw) {
    const query = authName ?
        JSON.stringify({ base: settings, remote: remote, auth: { name: authName, psw: authPsw } })
        : JSON.stringify({ base: settings, remote: remote });

    return fetchApi('/pull', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: query
    });
}

export function push(settings, remote, authName, authPsw) {
    const query = authName ?
        JSON.stringify({ base: settings, remote: remote, auth: { name: authName, psw: authPsw } })
        : JSON.stringify({ base: settings, remote: remote });

    return fetchApi('/push', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: query
    });
}

export function removeBranch(user, repo, branch) {
    const query = JSON.stringify({ repo: repo, user: user, branch: branch });

    return fetchApi('/branches', {
        method: 'DELETE',
        headers: {
            'Content-Type': 'application/json'
        },
        body: query
    });
}

export function removeRepo(user, repo) {
    const query = JSON.stringify({ repo: repo, user: user });

    return fetchApi('/repositories', {
        method: 'DELETE',
        headers: {
            'Content-Type': 'application/json'
        },
        body: query
    });
}

export function merge(settings, theirs) {
    const query = JSON.stringify({ base: settings, theirs: theirs })

    return fetchApi('/merge', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: query
    });
}

export function abortMerge(settings) {
    const query = JSON.stringify({ base: settings })

    return fetchApi('/merge/abort', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: query
    });
}
