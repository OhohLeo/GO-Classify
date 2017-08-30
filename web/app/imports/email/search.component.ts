import { Component, Input } from '@angular/core'
import { Search } from './email'

@Component({
    selector: 'email-search',
    templateUrl: './search.component.html'
})

export class EmailSearchComponent {
	@Input() search : Search

	onSince(value: string) {
		this.search.since = value
	}

	onBefore(value: string) {
		this.search.before = value
	}

	onText(event) {
		this.search.text = event.list
    }
}
