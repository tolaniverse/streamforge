
import { check, sleep } from 'k6';
import exec from 'k6/execution';
import http from 'k6/http';

export const options = {
	vus: 10,
	iterations: 100,
};

const BASE_URL = 'http://localhost:3000';

export default function() {
	const iteration = exec.scenario.iterationInTest;
	const eventId = `evt-${String(iteration + 1).padStart(3, '0')}`;

	const payload = JSON.stringify({
		event_type: 'payment.completed',
		event_id: eventId,
		payload: {
			transaction_id: 'txn-123',
			amount: 4500,
		},
	});

	const res = http.post(`${BASE_URL}/v1/events`, payload, {
		headers: {
			'Content-Type': 'application/json',
		},
	});

	check(res, {
		'status is 2xx': (r) => r.status >= 200 && r.status < 300,
	});

	console.log(`event_id=${eventId} status=${res.status}`);

	sleep(0.1);
}
