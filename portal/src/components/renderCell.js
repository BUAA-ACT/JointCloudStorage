export default {
  name: "render-cell",
  functional: true,
  props: {
    render: Function
  },
  render: (h, ctx) => ctx.props.render(h, ctx.data.attrs)
};
