import { Component, Input, OnInit } from '@angular/core'
import { Search } from './imap'

declare var jQuery: any;

@Component({
    selector: 'imap-search',
    templateUrl: './search.component.html'
})

export class ImapSearchComponent implements OnInit {
    @Input() search: Search

    ngOnInit() {
        jQuery('.datepicker').pickadate({
            selectMonths: true,
            selectYears: 15,
            today: 'Today',
            clear: 'Clear',
            close: 'Ok',
            closeOnSelect: false
        });
    }

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
