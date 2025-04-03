import { PUBLIC_API_URL } from '$env/static/public';

const categoriesURL = `${PUBLIC_API_URL}/categories`;

export async function fetchCategories(token: string) {
	try {
		const res = await fetch(categoriesURL, {
			headers: {
				Authorization: `Bearer ${token}`
			}
		});
		const data = await res.json();
		return data;
	} catch (error) {
		console.error('Failed to fetch experiment state:', error);
	}
	return [];
}
