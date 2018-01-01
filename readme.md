
## mtmapprune

A simple tool to prune minetests' map.sqlite file and delete all
blocks outside a certain range.

This program is licensed LGPL-2.1, with MIT components. It includes
some code from minetest verbatim to convert blockpos-nodepos.

## usage

`mtmapprune map.sqlite max_x [max_y [max_z [min_x min_y min_z]]]`

If omitted, `max_y` and `max_z` default to the value of the `max_x`
limit. The `min_*` limits will default to the negative values of the
`max_*` limits. The limits are node positions, not "block positions".

Example: to prune all nodes with x and z > 1000, and x and z < -1000
and y > 200, and y < 200, use:

`mtmapprune map.sqlite 1000 200`

