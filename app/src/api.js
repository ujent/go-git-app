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

export function switchUser(user) {
    const query = JSON.stringify({ name: user });

    return fetchApi('/users/switch', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: query
    });
}

export function getUsers() {
    return fetchApi('/users', {
    });
}

/* export function getBenefits(rq, user) {
    const query = queryString.stringify(rq);

    return fetchApi('/benefits?' + query, {
      headers: {
        Authorization: 'Bearer ' + user
      }
    });
  }

  export function stopFreeco(user) {
    return fetchApi('/stop', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: 'Bearer ' + user
      }
    });
  } */