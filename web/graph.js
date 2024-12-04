import * as d3 from "https://cdn.skypack.dev/d3@7";
import { uid } from "./utils.js";

class Graph {

	d3;
	svg;
	simulation;
	width;
	height;
	vertices = [];
	edges = [];

	constructor(parentDomNode) {
		const { width, height } = parentDomNode.getBoundingClientRect();
		this.width = width;
		this.height = height;
		this.svg = d3.create("svg").attr("width", width).attr("height", height);
		this.d3 = d3;
		//	buckets for nodes and links
		this.svg.append('g').classed('links', true);
		this.svg.append('g').classed('nodes', true);
		parentDomNode.append(this.svg.node());
		this.simulation = d3.forceSimulation();
	}

	runSimulation() {
		const s = this.simulation;
		const width = this.width;
		const height = this.height;
		const link = this.svg.selectAll('.links > line');
		const node = this.svg.selectAll('.nodes > g');
		//	reset before setting up our promise so we don't resolve too soon
		s.alpha(1).restart();
		return new Promise((resolve, reject) => {
			const tick = () => {
				link.attr("x1", d => d.source.x )
					.attr("y1", d => d.source.y )
					.attr("x2", d => d.target.x )
					.attr("y2", d => d.target.y );
				node.attr("transform", (d) => {
					return "translate(" + d.x + "," + d.y + ")";
				})
			};
			const strengthProportionalToCentrality = d => {
				const gamma = d.centrality;
				return gamma / 3333;
			};
			s.force("link", d3.forceLink().links(this.edges).strength(0.05).distance(85));
			s.force("charge", d3.forceManyBody().strength(-75).distanceMax(220));
			s.alphaMin(0.03);
			s.alphaDecay(0.3);
			s.force("collide", d3.forceCollide().radius(22));
			s.force('forceX', d3.forceX().x(width/2).strength(strengthProportionalToCentrality));
			s.force('forceY', d3.forceY().y(height/2).strength(strengthProportionalToCentrality));
			s.nodes(this.vertices).on("tick", tick);
			s.on("end", resolve);
		});
	}

	randomPointOnPlane() {
		//	margin ensures we're not spawning nodes along the edges
		const margin = 100;
		const x = Math.floor( Math.random()* (this.width-margin*2) + margin);
		const y = Math.floor( Math.random()* (this.height-margin*2)+ margin);
		return { x,y };
	}

	toVertex(obj) {
		/**
		 * @description:	takes an object, adds necessary properties
		 * @returns:		that new object
		 * @note:			the main purpose of this method should be to ensure proper structure
		 * 					either by adding necessary properties, or rejecting bad input
		 * 					by throwing an error. Guarantees about data structure go here
		 */
		const {x,y} = this.randomPointOnPlane();

		//	if the obj you're passing in has these props, they will be honored
		const defaultProps = {
			"_id": uid()
		};
		defaultProps.x = x;
		defaultProps.y = y;
		defaultProps.centrality = 1;
		return Object.assign(defaultProps, obj);
	}

	toEdge(linkObj) {
		if ("source" in linkObj && "target" in linkObj) {
			if (!("_id" in linkObj)) {
				linkObj._id = uid();
			}
		} else {
			throw new Error("source or link doesn't exist");
		}
		return linkObj;
	}

	vertexExists(vertex) {
		const nRows = this.vertices.filter(v => {
			return ("_id" in v) && (v._id === vertex._id);
		}).length;
		return !!(nRows);
	}

	onMouseOver(ev, d){
		//	no op
	}
	onMouseOut(ev, d){
		//	no op
	}
	onClick(){
		//	no op
	}

	async addVertex(vertex) {

		if (this.vertexExists(vertex)) {
			throw new Error("vertex exists");
		}
		this.vertices.push(vertex);

		//	SVG group that holds circle.node and text.label
		let nodesContainer = this.svg.select('.nodes');
		let node = nodesContainer.selectAll(".nodes > g")
			.data(this.vertices, d => d._id)
			.enter().append("g")
			.classed("node", true)
			.classed('tag', d => d.type==="tag")
			.classed('article', d => d.type==="article")
			.attr('x',d => d.x)
			.attr('y',d => d.y)
			.attr("transform",d => `translate(${d.x},${d.y})`)
			.on('mouseover', (ev, d) => {
				this.onMouseOver(ev, d);
			})
			.on('mouseout', (ev, d) => {
				this.onMouseOut(ev ,d);
			})
			.on('click', (ev, d) => {
				this.onClick(ev, d);
			})

		//	circle
		let circle = node.append("circle");
		//	drag & drop functionality
		const dragstarted = (ev) => {
			const d = ev.subject;
			if (!ev.active) {
				this.simulation.alphaTarget(0.3).restart();
			}
			d.fx = d.x;
			d.fy = d.y;
		};
		const dragged = (ev) => {
			const d = ev.subject;
			d.fx = ev.x;
			d.fy = ev.y;
		};
		const dragended = (ev) => {
			const d = ev.subject;
			if (!ev.active) {
				this.simulation.alphaTarget(0);
			}
			d.fx = null;
			d.fy = null;
		};
		const drag_handler = d3.drag()
			.on("start", dragstarted)
			.on("drag", dragged)
			.on("end", dragended);
		drag_handler(node);

		node.append("text")
			.text(d => d.title)
			.attr('x', 17)
			.attr('y', 15);
		
		return await this.runSimulation();
	}
	async addEdge(edge) {
		this.edges.push(edge);

		this.svg.selectAll(".links")
		.selectAll("line")
		.data(this.edges)
		.join(
			enter => enter.append("line").classed("link", true),
			update => update,
			exit => exit.remove()
		);

		//	recalculate centrality for every vertex
		this.vertices.forEach(v => {
			v.centrality=0;
		});
		this.edges.forEach(e => {
			e.source.centrality++;
			e.target.centrality++;
		});


		return await this.runSimulation();

	}

}

export default Graph;