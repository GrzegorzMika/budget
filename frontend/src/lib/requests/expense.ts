import { PUBLIC_API_URL } from '$env/static/public';

const expenseURL = `${PUBLIC_API_URL}/expenses`;

export async function createExpenses(data: FormData) {
	const request = Object.fromEntries(data);
	const { jwt, ...rest } = request;
	try {
		await fetch(expenseURL, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
				Authorization: `Bearer ${jwt}`
			},
			body: JSON.stringify(rest)
		});
	} catch (error) {
		console.error('Failed to create expense:', error);
	}
}
