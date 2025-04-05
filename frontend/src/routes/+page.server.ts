import { createExpenses } from '$lib';
import type { Actions } from './$types';

export const actions = {
	default: async ({ request }) => {
		const data = await request.formData();
		await createExpenses(data);
	}
} satisfies Actions;
