# Goal
**Find a node ID arrangement algorithm, so that the E(routing_distance) is lower than average distance between two Random points (discrete uniform distribution) in a N*N grid.**


## Average Manhattan Distance Between Two Random Points

## Manhattan Distance
* Definition: The distance between two points measured along axes at right angles. In a plane with $p_1(x_1, y_1)$ and $p_2(x_2, y_2)$, it is $|x1 - x2| + |y1 - y2|$.


## Math
$x_1, y_1, x_2, y_2$ are follow a discrete uniform distribution from 0 to N-1.

Let $Y = |x1 - x2| + |y1 - y2|$,

$E(Y) = E(|x1 - x2|) + E(|y1 - y2|) = 2E(|x1 - x2|)$


Define function $g (x_1,x_2) = |x_1 - x_2| = \begin{cases} x_1 - x_2 & \textrm{if} \quad{} x_1 \geq x_2\\ x_2 - x_1 & \textrm{if} \quad{} x_2 \geq x_1\end{cases}$

Two points are picked independently, thus the joint probability mass function is product of the PMFs