(function () {
	'use strict';

	// Drop the one-shot ?saved=1 flag from the URL so a refresh doesn't
	// re-show the confirmation.
	if (window.location.search.indexOf('saved=1') > -1 && window.history.replaceState) {
		window.history.replaceState(null, '', '/?tab=add');
	}

	// Paper/ink mood toggle; the initial value is applied by an inline
	// head script before first paint.
	var mood = document.getElementById('mood-toggle');
	if (mood) {
		mood.addEventListener('click', function () {
			var root = document.documentElement;
			var next = root.dataset.theme === 'ink' ? 'paper' : 'ink';
			root.dataset.theme = next;
			try { localStorage.setItem('mood', next); } catch (e) { /* private mode */ }
		});
	}

	// Searchable category picker on the add-expense form.
	var picker = document.getElementById('catpicker');
	if (picker) {
		var btn = picker.querySelector('.catpicker-btn');
		var pop = picker.querySelector('.catpicker-pop');
		var search = picker.querySelector('.catpicker-search input');
		var hidden = picker.querySelector('input[name="category"]');
		var nameEl = picker.querySelector('.catpicker-name');
		var dotEl = btn.querySelector('.dot');
		var empty = picker.querySelector('.catpicker-empty');
		var opts = Array.prototype.slice.call(picker.querySelectorAll('.catpicker-opt'));

		var close = function () { pop.hidden = true; };
		var open = function () {
			pop.hidden = false;
			search.value = '';
			filter('');
			search.focus();
		};
		var filter = function (q) {
			q = q.trim().toLowerCase();
			var visible = 0;
			opts.forEach(function (o) {
				var hit = o.dataset.name.toLowerCase().indexOf(q) > -1;
				o.hidden = !hit;
				if (hit) visible++;
			});
			empty.hidden = visible > 0;
		};

		btn.addEventListener('click', function () { pop.hidden ? open() : close(); });
		search.addEventListener('input', function () { filter(search.value); });
		opts.forEach(function (o) {
			o.addEventListener('click', function () {
				hidden.value = o.dataset.name;
				nameEl.textContent = o.dataset.name;
				nameEl.className = 'catpicker-name';
				dotEl.className = 'dot bg-' + o.dataset.color;
				opts.forEach(function (x) { x.classList.toggle('selected', x === o); });
				close();
			});
		});

		// No category picked yet: open the picker instead of submitting.
		var form = picker.closest('form');
		if (form) {
			form.addEventListener('submit', function (e) {
				if (!hidden.value) {
					e.preventDefault();
					open();
				}
			});
		}
		document.addEventListener('click', function (e) {
			if (!picker.contains(e.target)) close();
		});
		document.addEventListener('keydown', function (e) {
			if (e.key === 'Escape') close();
		});
	}
})();
