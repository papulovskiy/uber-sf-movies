#!/usr/bin/perl

use strict;
use warnings;

use JSON;
use LWP::UserAgent;
use Text::CSV;

die("Please provide Google Maps API key as a first argument") unless $ARGV[0];
my $key = $ARGV[0];

my $csv = Text::CSV->new({ binary => 1, eol    => $/, }) or die "Cannot use CSV: " . Text::CSV->error_diag();

open(my $input, "<:encoding(utf8)", "./downloaded_list.csv") or die "Cannot open downloaded_list.csv: $!";
open(my $output, ">:encoding(utf8)", "./coded_list.csv")     or die "Cannot open coded_list.csv: $!";

# Pass the header straight to the output file
# with additional fields for coordinates
my $header = $csv->getline($input);
push(@$header, "latitude", "longitude");
$csv->print($output, $header);

# Process data
my $ua = LWP::UserAgent->new(ssl_opts => { verify_hostname => 0 });
my $counter = 0;
while(my $row = $csv->getline($input)) {
    my $location = $row->[2];
    $location .= ", San Francisco";

    # Geocoding magic
    my $geo = $ua->get("https://maps.googleapis.com/maps/api/geocode/json?address=$location&key=$key&region=us");
    my $result = {};
    if($geo->is_success()) {
        eval {
            $result = from_json($geo->decoded_content);
            1;
        };
    }
    if(defined($result->{results}->[0]->{geometry}->{location})) {
        push(@$row, $result->{results}->[0]->{geometry}->{location}->{lat});
        push(@$row, $result->{results}->[0]->{geometry}->{location}->{lng});
    } else {
        push(@$row, "", ""); # to keep the same number of columns
    }

    $csv->print($output, $row);
    $counter++;
}

print "$counter\n"; # just to know how many rows passed


close($input);
close($output);