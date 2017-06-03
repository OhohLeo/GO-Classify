import { Component, Input, OnInit } from '@angular/core';
import { Movie } from './movie'

@Component({
    selector: 'simple-movie',
    templateUrl: './simple.component.html'
})

export class SimpleMovieComponent implements OnInit {
    @Input() movie: Movie

    ngOnInit() {
    }
}
