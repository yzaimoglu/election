// @ts-nocheck
const base = 'http://localhost:84/v1';

async function send({ method, path, data, token }) {
	const opts = { method, headers: {} };

	if (data) {
		opts.headers['Content-Type'] = 'application/json';
		opts.body = JSON.stringify(data);
	}

	return fetch(`${base}/${path}`, opts)
		.then((r) => r.text())
		.then((json) => {
			try {
				var resParsed = JSON.parse(json);

				if (resParsed?.status === 'error') {
					console.log(`API response error from ${base}/${path}: ${json}`);
				}

				return resParsed;
			} catch (err) {
				return json;
			}
		});
}

export function get(path) {
	return send({ method: 'GET', path });
}

export function del(path) {
	return send({ method: 'DELETE', path });
}

export function post(path, data) {
	return send({ method: 'POST', path, data });
}

export function put(path, data) {
	return send({ method: 'PUT', path, data });
}