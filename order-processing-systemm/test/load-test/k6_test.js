import http from 'k6/http';

import { sleep } from 'k6';

export const options = {
    vus: 10,
    duration: '10s',
};

export default function () {

    const payload = JSON.stringify({
        user_id: 1,
        product_id: 101,
        quantity: 2,
        amount: 500
    });

    const params = {
        headers: {
            'Content-Type': 'application/json',
        },
    };

    http.post(
        'http://localhost:8000/orders',
        payload,
        params,
    );

    sleep(1);
}