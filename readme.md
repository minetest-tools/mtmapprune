
## mtmapprune

A simple tool to prune minetests' map.sqlite file and delete all
blocks outside a certain range.

This program is licensed LGPL-2.1, with MIT components. It includes
some code from minetest verbatim to convert blockpos-nodepos.

## usage

`mtmapprune map.sqlite x_limit [y_limit] [z_limit]`

If omitted, `y_limit` and `z_limit` default to the value of the
`x_limit`. The limits are node positions, not "block positions".

Example: to prune all nodes with x and z > 1000, and x and z < -1000
and y > 200, and y < 200, use:

`mtmapprune map.sqlute 1000 200`

