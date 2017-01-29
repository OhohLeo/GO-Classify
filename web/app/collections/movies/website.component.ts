import { Component, NgZone, Input, OnInit, OnDestroy } from '@angular/core';
import { Movie } from './movie'

@Component({
	selector: 'website-movie',
	templateUrl: './website.component.html'
})

export class WebsiteMovieComponent implements OnInit, OnDestroy {

	@Input() movie: Movie

	private needDetails: boolean = false

    constructor(private zone: NgZone) {
	}

    ngOnInit() {
    }

    ngOnDestroy() {
    }

	getDetails() {
		this.zone.run(() => {
			this.needDetails = !this.needDetails;
		})
	}
}
