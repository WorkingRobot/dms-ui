// Cryptographically-strong password generator for the add/change dialogs.
// Avoids ambiguous chars (O/0, l/1, I) and shell-hostile quotes/backticks so
// generated passwords are safe to paste anywhere.
const CHARS = 'ABCDEFGHJKMNPQRSTUVWXYZabcdefghijkmnpqrstuvwxyz23456789!@#$%^&*-_=+?';

export function generatePassword(len = 20) {
	const buf = new Uint32Array(len);
	crypto.getRandomValues(buf);
	let out = '';
	for (let i = 0; i < len; i++) out += CHARS[buf[i] % CHARS.length];
	return out;
}

// Rough strength estimate (0-4) for the meter on password fields.
export function strength(pw) {
	if (!pw) return 0;
	let s = 0;
	if (pw.length >= 12) s++;
	if (pw.length >= 16) s++;
	if (/[a-z]/.test(pw) && /[A-Z]/.test(pw) && /\d/.test(pw)) s++;
	if (/[^A-Za-z0-9]/.test(pw)) s++;
	return Math.min(s, 4);
}
