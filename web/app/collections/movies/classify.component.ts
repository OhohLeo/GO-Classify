import { Component, NgZone, Input, OnInit, OnDestroy } from '@angular/core'
import { Movie } from './movie'

@Component({
	selector: 'classify-movie',
	templateUrl: './classify.component.html'
})

export class ClassifyMovieComponent implements OnInit, OnDestroy {

	@Input() movie: Movie

    ngOnInit() {
    }

    ngOnDestroy() {
    }
}
