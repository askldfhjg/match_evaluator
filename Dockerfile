FROM alpine
ADD match_evaluator /match_evaluator
ENTRYPOINT [ "/match_evaluator" ]
