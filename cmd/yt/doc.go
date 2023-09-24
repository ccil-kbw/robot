/*
package main

TODO: This doc file serves as notes on how to handle the automated-youtube uploader. Please delete and write proper docs when the feature is delivered

Issue:
Currently we have videos that are recorded every day from the great lessons @ the Masjid Khaled Ben Walid,
the Videos are recorded using OBS Software and are saved as .mp4 and .mkv in a predictable path based on the date (YYYY-MM-dd/.*\.(mp4|mkv))

The Videos are then picked by the team in the Cloud Folder, trimmed and uploaded manually to YouTube.

The issue is the backlog is extremely big and having to trim these manually and upload them will be a pain for the Operation team.

Solution:
Upload the YouTube Videos automatically in a private playlist (flag the videos as private too), then the only thing left is to Trim the video, give a title and make it public

The advantage of such a solution is that the heavy-lifting of looking for the video in the Cloud folder and uploading it to YouTube is done automatically,
the manual part that is left is (hopefully) much simpler, we only need to trim the video which can be done on https://studio.youtube.com from any computer (no need of special software. Might even work from tablets or phone to-be-verified)
*/
package main
