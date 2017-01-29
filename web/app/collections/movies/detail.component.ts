import { Component, NgZone, Input, OnInit, OnDestroy } from '@angular/core';
import { Movie } from './movie'

@Component({
	selector: 'detail-movie',
	templateUrl: './detail.component.html'
})

export class DetailMovieComponent implements OnInit, OnDestroy {

	@Input() movie: Movie

    ngOnInit() {
		console.log(this.movie)
    }

    ngOnDestroy() {
    }
}
