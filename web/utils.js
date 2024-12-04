const waitFor = (ms) => {
	return new Promise((resolve, reject) => {
		if (typeof ms !== "number") {
			reject("expected ms to be a number. Got " + typeof ms);
		}
		const resolveTrue = () => resolve({"waitedFor": ms, "unit": "milliseconds"});
		setTimeout(resolveTrue, ms);
	});
};

const uid = () => {
	//	a cheap but reliably unique ID
	return Math.floor((Math.random() * Number.MAX_SAFE_INTEGER)).toString(16).toUpperCase();
}

export { uid, waitFor };
