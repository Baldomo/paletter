/*
Package paletter provides methods to extract a color palette from an image through clustering (k-Means).
Given an image, paletter stores all colors in the L*a*b* (CIELAB) color space, then finds the centers
of clusters created from said colors. The clusters' center represent the most prominent color in the
original image.

The CIELAB color space was chosen both for the simplicity of its representation and it being device-indipendent
(it also happens to be copyright and license-free).

The subcommand cmd/paletter makes use of this package to generate images laid out to contain both the
original image and a palette with the prominent colors.
*/

package paletter
